package main

import (
	"log"
	"net/http"
)

func teapotHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/art/favicon.ico")
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./assets/docs")))
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/teapot", teapotHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
