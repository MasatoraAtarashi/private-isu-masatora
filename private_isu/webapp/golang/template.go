package main

import (
	"text/template"
)

var (
	tplCache = make(map[string]*template.Template, 6)
)

const (
	templateKeyGetIndex = "layout.html"
	templateKeyGetPost  = "posts.html"
)

func parseTemplates() {
	fmap := template.FuncMap{
		"imageURL": imageURL,
	}

	tplCache[templateKeyGetIndex] = template.Must(template.New(templateKeyGetIndex).Funcs(fmap).ParseFiles(
		getTemplPath("layout.html"),
		getTemplPath("index.html"),
		getTemplPath("posts.html"),
		getTemplPath("post.html"),
	))
}
