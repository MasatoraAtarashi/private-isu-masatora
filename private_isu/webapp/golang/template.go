package main

import (
	"text/template"
)

var (
	tplCache = make(map[string]*template.Template, 7)
)

const (
	templateKeyGetIndex       = "getIndex"
	templateKeyGetAccountName = "getAccountName"
	templateKeyGetAdminBanned = "getAdminBanned"
	templateKeyGetLogin       = "getLogin"
	templateKeyGetRegister    = "getRegister"
	templateKeyGetPostsID     = "getPostsID"
)

func parseTemplates() {
	fmap := template.FuncMap{
		"imageURL": imageURL,
	}

	tplCache[templateKeyGetIndex] = template.Must(template.New("layout.html").Funcs(fmap).ParseFiles(
		getTemplPath("layout.html"),
		getTemplPath("index.html"),
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

	tplCache[templateKeyGetPostsID] = template.Must(template.New("layout.html").Funcs(fmap).ParseFiles(
		getTemplPath("layout.html"),
		getTemplPath("post_id.html"),
		getTemplPath("post.html"),
	))
}
