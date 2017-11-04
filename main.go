package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("view/*.html"))
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/about", about)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// the home page
func home(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// the about page
func about(w http.ResponseWriter, req *http.Request) {
	//io.WriteString(w, "about")
	err := tpl.ExecuteTemplate(w, "about.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}
