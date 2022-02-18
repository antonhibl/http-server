package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/hello-hypertext", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello friend."))
	})
	http.ListenAndServe(":8080", nil)
}
