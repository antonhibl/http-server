package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type BlogPost struct {
	Title       string `json:"title"`
	Timestamp   string `json:"timestamp"`
	Main        string `json:"main"`
	ContentInfo string `json:"content_info"`
}

var blogTemplate = template.Must(template.ParseFiles("./blogtemplate.html"))

func blogHandler(w http.ResponseWriter, r *http.Request) {
	blogstr := r.URL.Path[len("/blog/"):] + ".json"

	f, err := os.Open(blogstr)
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
	if err := blogTemplate.Execute(w, post); err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/blog/", blogHandler)
	// To serve a directory on disk (/path/to/assets/on/my/computer)
	// under an alternate URL path (/assets/), use StripPrefix to
	// modify the request URL's path before the FileServer sees it:
	http.Handle("/assets/", http.StripPrefix("/assets/",
		http.FileServer(http.Dir("/Users/cthulhu/l4p1s/languages/ECMAScript/learning/golangWeb/"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
