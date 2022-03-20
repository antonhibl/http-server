package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type BlogPost struct {
	Title      string   `json:"title"`
	Timestamp  string   `json:"timestamp"`
	Main       []string `json:"main"`
	ParsedMain template.HTML
}

var blogTemplate = template.Must(template.ParseFiles("./assets/docs/blogtemplate.html"))

func blogHandler(w http.ResponseWriter, r *http.Request) {
	blogstr := r.URL.Path[len("/blog/"):] + ".json"

	f, err := os.Open("db/" + blogstr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer f.Close()

	var post BlogPost
	if err := json.NewDecoder(f).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post.ParsedMain = template.HTML(strings.Join(post.Main, ""))

	if err := blogTemplate.Execute(w, post); err != nil {
		log.Println(err)
	}
}

func teapotHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("<html><h1><a href='https://datatracker.ietf.org/doc/html/rfc2324/'>HTCPTP</h1><img src='https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Ftaooftea.com%2Fwp-content%2Fuploads%2F2015%2F12%2Fyixing-dark-brown-small.jpg&f=1&nofb=1' alt='Im a teapot'><html>"))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/art/favicon.ico")
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("/Users/cthulhu/l4p1s/languages/Go/net/http-server/assets/docs")))
	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/teapot", teapotHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
