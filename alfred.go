package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
)

func githubnotify(w http.ResponseWriter, r *http.Request) {
	spew.Dump(r)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		spew.Dump(err)
	}
	fmt.Println(string(data))
}

func main() {
	http.HandleFunc("/githubnotify/", githubnotify)
	http.ListenAndServe(":8080", nil)
}
