package main

import (
	"text/template"
)

var (
	tplCache = make(map[string]*template.Template, 6)
)

const (
	templateKeyGetIndex       = "getIndex"
	templateKeyGetAccountName = "getAccountName"
	templateKeyGetAdminBanned = "getAdminBanned"
	templateKeyGetLogin       = "getLogin"
	templateKeyGetRegister    = "getRegister"
	templateKeyGetPost        = "getPost"
	templateKeyGetPostsID     = "getPostsID"
)

func parseTemplates() {
	fmap := template.FuncMap{
		"imageURL": imageURL,
	}

	tplCache[templateKeyGetIndex] = template.Must(template.New("layout.html").Funcs(fmap).ParseFiles(
		getTemplPath("layout.html"),
		getTemplPath("index.html"),
		getTemplPath("posts.html"),
		getTemplPath("post.html"),
	))

	tplCache[templateKeyGetAccountName] = template.Must(template.New("layout.html").Funcs(fmap).ParseFiles(
		getTemplPath("layout.html"),
		getTemplPath("user.html"),
		getTemplPath("posts.html"),
		getTemplPath("post.html"),
	))

	tplCache[templateKeyGetAdminBanned] = template.Must(template.ParseFiles(
		getTemplPath("layout.html"),
		getTemplPath("banned.html")),
	)

	tplCache[templateKeyGetLogin] = template.Must(template.ParseFiles(
		getTemplPath("layout.html"),
		getTemplPath("login.html")),
	)

	tplCache[templateKeyGetRegister] = template.Must(template.ParseFiles(
		getTemplPath("layout.html"),
		getTemplPath("register.html")),
	)

	tplCache[templateKeyGetPost] = template.Must(template.New("posts.html").Funcs(fmap).ParseFiles(
		getTemplPath("posts.html"),
		getTemplPath("post.html"),
	))

	tplCache[templateKeyGetPostsID] = template.Must(template.New("layout.html").Funcs(fmap).ParseFiles(
		getTemplPath("layout.html"),
		getTemplPath("post_id.html"),
		getTemplPath("post.html"),
	))
}
