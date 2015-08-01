package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/http"
)

func githubnotify(w http.ResponseWriter, r *http.Request) {
	spew.Print(r)
}

func main() {
	http.HandleFunc("/githubnotify/", githubnotify)
	http.ListenAndServe(":8080", nil)
}
