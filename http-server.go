package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

// Define basic blog post structure
type BlogPost struct {
	Title      string   `json:"title"`
	Timestamp  string   `json:"timestamp"`
	Main       []string `json:"main"`
	ParsedMain template.HTML
}

// require my blog HTML template for a template parser
var blogTemplate = template.Must(template.ParseFiles("./assets/docs/blogtemplate.html"))

// Handler to dynamically build blog post webpages from my JSON database
func blogHandler(w http.ResponseWriter, r *http.Request) {
	// locate JSON file
	blogstr := r.URL.Path[len("/blog/"):] + ".json"

	// open json file
	f, err := os.Open("db/" + blogstr)
	// if no file found
	if err != nil {
		// return error status
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// defer the closing of this json file until the page is done loading content
	defer f.Close()

	// define a blog post object
	var post BlogPost
	// decode the JSON data into the object
	if err := json.NewDecoder(f).Decode(&post); err != nil {
		// if an error occurs return HTTP status
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// parse the post's object data into the template
	post.ParsedMain = template.HTML(strings.Join(post.Main, " "))

	// execute and serve the template
	if err := blogTemplate.Execute(w, post); err != nil {
		// if an error occurs return status
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// teapot joke handler :: redirects to HTCPCP protocol
func teapotHandler(w http.ResponseWriter, r *http.Request) {
	// return teapot state
	w.WriteHeader(http.StatusTeapot)
	// serve teapot image with linking
	w.Write([]byte("<html><h1><a href='https://datatracker.ietf.org/doc/html/rfc2324/'>HTCPTP</h1><img src='https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Ftaooftea.com%2Fwp-content%2Fuploads%2F2015%2F12%2Fyixing-dark-brown-small.jpg&f=1&nofb=1' alt='Im a teapot'><html>"))
}

// favicon handler
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	// serve the favicon to client
	http.ServeFile(w, r, "./assets/art/favicon.ico")
}

// main program
func main() {
	// define the file server and its root dir
	http.Handle("/", http.FileServer(http.Dir("/Users/cthulhu/l4p1s/languages/Go/net/http-server/assets/docs")))
	// initialize blog handler
	http.HandleFunc("/blog/", blogHandler)
	// initialize favicon handler
	http.HandleFunc("/favicon.ico", faviconHandler)
	// initialize teapot handler
	http.HandleFunc("/teapot", teapotHandler)
	// listen and serve on port 8080, log and exit safely if error occurs.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
