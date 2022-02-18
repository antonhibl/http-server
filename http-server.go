package main

import (
	"log"
	"net/http"
)

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/art/favicon.ico")
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./assets/docs")))
	http.HandleFunc("/favicon.ico", faviconHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
