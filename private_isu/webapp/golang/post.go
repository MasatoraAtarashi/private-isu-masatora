package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/go-chi/chi/v5"
)

type Post struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	Imgdata      []byte    `db:"imgdata"`
	Body         string    `db:"body"`
	Mime         string    `db:"mime"`
	CreatedAt    time.Time `db:"created_at"`
	CommentCount int
	Comments     []Comment
	User         User
	CSRFToken    string
}

func makePosts(results []Post, csrfToken string, allComments bool) ([]Post, error) {
	var posts []Post

	var countKeys []string
	var commentKeys []string
	for _, p := range results {
		countKeys = append(countKeys, toCommentCountCacheKey(p.ID))
		commentKeys = append(commentKeys, toCommentCacheKey(p.ID, allComments))
	}

	// コメント数をmemcachedから取得
	cachedCommentCountMap, err := memcacheClient.GetMulti(countKeys)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		return nil, err
	}

	cachedCommentMap, err := memcacheClient.GetMulti(commentKeys)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		return nil, err
	}

	for _, p := range results {
		// キャッシュが存在したらそれを使う
		if cachedCommentCount := cachedCommentCountMap[toCommentCountCacheKey(p.ID)]; cachedCommentCount != nil {
			// どうやるんや
			countStrBrackets := string(cachedCommentCount.Value)
			countStrBrackets = strings.Replace(countStrBrackets, "[", "", -1)
			countStr := strings.Replace(countStrBrackets, "]", "", -1)
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return nil, err
			}

			p.CommentCount = count
		} else {
			// 存在しなかったらDBに問い合せる
			err = db.Get(&p.CommentCount, "SELECT COUNT(*) AS `count` FROM `comments` WHERE `post_id` = ?", p.ID)
			if err != nil {
				return nil, err
			}

			countStr := strconv.Itoa(p.CommentCount)

			// キャッシュをセット
			err = memcacheClient.Set(&memcache.Item{
				Key:   fmt.Sprintf("comments.%d.count", p.ID),
				Value: []byte(countStr),
			})
			if err != nil {
				return nil, err
			}
		}

		if cachedComments := cachedCommentMap[toCommentCacheKey(p.ID, allComments)]; cachedComments != nil {
			var comments []Comment
			err := json.Unmarshal(cachedComments.Value, &comments)
			if err != nil {
				return nil, err
			}
			p.Comments = comments
		} else {
			query := "SELECT * FROM `comments` WHERE `post_id` = ? ORDER BY `created_at` DESC"
			if !allComments {
				query += " LIMIT 3"
			}
			var comments []Comment
			err = db.Select(&comments, query, p.ID)
			if err != nil {
				return nil, err
			}

			for i := 0; i < len(comments); i++ {
				err := db.Get(&comments[i].User, "SELECT * FROM `users` WHERE `id` = ?", comments[i].UserID)
				if err != nil {
					return nil, err
				}
			}

			// reverse
			for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
				comments[i], comments[j] = comments[j], comments[i]
			}

			p.Comments = comments

			err := setCommentsCache(comments, p.ID, allComments)
			if err != nil {
				return nil, err
			}
		}

		p.CSRFToken = csrfToken

		posts = append(posts, p)
		if len(posts) >= postsPerPage {
			break
		}
	}

	return posts, nil
}

func setCommentsCache(comments []Comment, pid int, allComments bool) error {
	bytes, err := json.Marshal(comments)
	if err != nil {
		return err
	}

	err = memcacheClient.Set(&memcache.Item{
		Key:        fmt.Sprintf("comments.%d.%t", pid, allComments),
		Value:      bytes,
		Flags:      0,
		Expiration: 0,
	})

	return err
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	maxCreatedAt := m.Get("max_created_at")
	if maxCreatedAt == "" {
		return
	}

	t, err := time.Parse(ISO8601Format, maxCreatedAt)
	if err != nil {
		log.Print(err)
		return
	}

	results := []Post{}
	err = db.Select(&results, "SELECT `id`, `user_id`, `body`, `mime`, `created_at` FROM `posts` WHERE `created_at` <= ? ORDER BY `created_at` DESC LIMIT ?", t.Format(ISO8601Format), postsPerPage)
	if err != nil {
		log.Print(err)
		return
	}

	posts, err := makePosts(results, getCSRFToken(r), false)
	if err != nil {
		log.Print(err)
		return
	}

	if len(posts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(getPostsHTMLString(posts)))
}

func getPostsID(w http.ResponseWriter, r *http.Request) {
	pidStr := chi.URLParam(r, "id")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	results := []Post{}
	err = db.Select(&results, "SELECT * FROM `posts` WHERE `id` = ?", pid)
	if err != nil {
		log.Print(err)
		return
	}

	posts, err := makePosts(results, getCSRFToken(r), true)
	if err != nil {
		log.Print(err)
		return
	}

	if len(posts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	p := posts[0]

	me := getSessionUser(r)

	tplCache[templateKeyGetPostsID].Execute(w, struct {
		Post Post
		Me   User
	}{p, me})
}

func postIndex(w http.ResponseWriter, r *http.Request) {
	me := getSessionUser(r)
	if !isLogin(me) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.FormValue("csrf_token") != getCSRFToken(r) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		session := getSession(r)
		session.Values["notice"] = "画像が必須です"
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	mime := ""
	ext := ""
	if file != nil {
		// 投稿のContent-Typeからファイルのタイプを決定する
		contentType := header.Header["Content-Type"][0]
		if strings.Contains(contentType, "jpeg") {
			mime = "image/jpeg"
			ext = "jpg"
		} else if strings.Contains(contentType, "png") {
			mime = "image/png"
			ext = "png"
		} else if strings.Contains(contentType, "gif") {
			mime = "image/gif"
			ext = "gif"
		} else {
			session := getSession(r)
			session.Values["notice"] = "投稿できる画像形式はjpgとpngとgifだけです"
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	filedata, err := io.ReadAll(file)
	if err != nil {
		log.Print(err)
		return
	}

	if len(filedata) > UploadLimit {
		session := getSession(r)
		session.Values["notice"] = "ファイルサイズが大きすぎます"
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	query := "INSERT INTO `posts` (`user_id`, `mime`, `imgdata`, `body`) VALUES (?,?,?,?)"
	result, err := db.Exec(
		query,
		me.ID,
		mime,
		[]byte(""),
		r.FormValue("body"),
	)
	if err != nil {
		log.Print(err)
		return
	}

	pid, err := result.LastInsertId()
	if err != nil {
		log.Print(err)
		return
	}

	err = writeImageFile(pid, ext, filedata)
	if err != nil {
		log.Print(err)
		return
	}

	http.Redirect(w, r, "/posts/"+strconv.FormatInt(pid, 10), http.StatusFound)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	pidStr := chi.URLParam(r, "id")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	post := Post{}
	err = db.Get(&post, "SELECT imgdata, mime FROM `posts` WHERE `id` = ?", pid)
	if err != nil {
		log.Print(err)
		return
	}

	ext := chi.URLParam(r, "ext")

	if ext == "jpg" && post.Mime == "image/jpeg" ||
		ext == "png" && post.Mime == "image/png" ||
		ext == "gif" && post.Mime == "image/gif" {
		w.Header().Set("Content-Type", post.Mime)
		_, err := w.Write(post.Imgdata)
		if err != nil {
			log.Print(err)
			return
		}
		return
	}

	err = writeImageFile(int64(pid), ext, post.Imgdata)
	if err != nil {
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func writeImageFile(pid int64, ext string, data []byte) error {
	// 画像ファイルをpublic/image/ディレクトリに書き込む
	imgFile, err := os.Create(fmt.Sprintf("%s/%s.%s", imgDir, strconv.FormatInt(pid, 10), ext))
	if err != nil {
		return err
	}
	defer imgFile.Close()

	writer := bufio.NewWriter(imgFile)
	_, err = writer.Write(data)
	return err
}
