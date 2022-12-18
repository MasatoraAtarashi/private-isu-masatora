package main

import (
	"errors"
	"log"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	gsm "github.com/bradleypeabody/gorilla-sessions-memcache"
)

func init() {
	memdAddr := os.Getenv("ISUCONP_MEMCACHED_ADDRESS")
	if memdAddr == "" {
		memdAddr = "localhost:11211"
	}
	memcacheClient = memcache.New(memdAddr)
	store = gsm.NewMemcacheStore(memcacheClient, "iscogram_", []byte("sendagaya"))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	err := toStaticImageFile()
	if err != nil {
		panic(err)
	}

	parseTemplates()
}

// MySQLに最初から入ってる画像を静的ファイル化する
func toStaticImageFile() error {
	var postIDs []int
	err := db.Get(&postIDs, "SELECT id FROM posts")
	if err != nil {
		log.Print(err)
		return err
	}

	for _, postID := range postIDs {
		post := Post{}
		err = db.Get(&post, "SELECT imgdata, mime FROM `posts` WHERE `id` = ?", postID)
		if err != nil {
			log.Print(err)
			return err
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
			return errors.New("unknown mime")
		}

		err = writeImageFile(int64(postID), ext, post.Imgdata)
	}

	return nil
}
