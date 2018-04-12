package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Data interface {
}

func TemplateHandler(path string) http.HandlerFunc {
	const indexPage = "/index.html"
	return func(w http.ResponseWriter, r *http.Request) {
		var data Data

		page := r.URL.Query().Get(":page")

		dir, _ := os.Getwd()
		thepath := filepath.Join(dir, path+page)
		log.Println("TemplateHandler, r.URL.Path -> ", r.URL.Path, thepath)

		info, err := os.Stat(thepath)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
		}
		if info.IsDir() {
			thepath = thepath + indexPage
		}
		t, err := template.ParseFiles(filepath.Join(thepath), filepath.Join(filepath.Join(dir, "static/site/startup-kit/"), "layout.html"))

		if err != nil {
			fmt.Println("TemplateHandler", err)
			fmt.Fprintln(w, err)
			return
		}
		t.Execute(w, data)

	}
}
