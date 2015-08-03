package main

import (
	"encoding/json"
	//"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
)

func githubnotify(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		spew.Dump(err)
		return
	}

	payload := GithubPayload{}
	err = json.Unmarshal(data, &payload)
	if err != nil {
		spew.Dump(err)
		return
	}

	spew.Dump(payload)
}

func main() {
	http.HandleFunc("/githubnotify/", githubnotify)
	http.ListenAndServe(":8080", nil)
}
