package main

import (
	"errors"
	"log"
	"net/http"
)

// MySQLに最初から入ってる画像を静的ファイル化する
func toStaticImageFile(w http.ResponseWriter, r *http.Request) {
	var posts []Post
	err := db.Get(&posts, "SELECT id FROM posts")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}

	log.Println("default images count", len(posts))

	for _, post := range posts {
		log.Println("to static file: ", post.ID)

		post := Post{}
		err = db.Get(&post, "SELECT imgdata, mime FROM `posts` WHERE `id` = ?", post.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		}

		var ext string
		switch post.Mime {
		case "image/jpeg":
			ext = "jpeg"
		case "image/png":
			ext = "png"
		case "image/gif":
			ext = "gif"
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(errors.New("unknown mime"))
			return
		}

		err = writeImageFile(int64(post.ID), ext, post.Imgdata)
	}

	w.WriteHeader(http.StatusOK)
}
