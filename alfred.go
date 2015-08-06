package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//	"strings"
)

var secretKey []byte

func verifyHMAC(w http.ResponseWriter, r *http.Request, body []byte) bool {
	sig := r.Header.Get("X-Hub-Signature")

	mac := hmac.New(sha1.New, secretKey)
	mac.Write(body)
	expectedMAC := "sha1=" + hex.EncodeToString(mac.Sum(nil))
	equal := hmac.Equal([]byte(sig), []byte(expectedMAC))

	if !equal {
		http.Error(w, "X-Hub-Signature Mismatch", http.StatusUnauthorized)
	}

	return equal
}

func githubnotify(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Header.Get("X-Github-Event") != "push" {
		return
	}

	if !verifyHMAC(w, r, body) {
		return
	}

	payload := GithubPayload{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := payload.Repository.Description
	url := payload.Repository.URL
	fmt.Println(name, url)

}

func main() {
	secretKey = []byte(os.Getenv("GITHUBSECRET"))
	if len(secretKey) == 0 {
		panic("GITHUBSECRET environment variable not set")
	}

	http.HandleFunc("/githubnotify/", githubnotify)
	http.ListenAndServe(":8080", nil)
}
