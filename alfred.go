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
	"strings"
)

func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

var secretKey []byte

func githubnotify(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event := r.Header.Get("X-Github-Event")
	sig := r.Header.Get("X-Hub-Signature")

	payload := GithubPayload{}
	err = json.Unmarshal(data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.HasPrefix(sig, "sha1=") {
		sigBytes, err := hex.DecodeString(sig[len("sha1="):])
		if err != nil {
			http.Error(w, "Invalid X-Hub-Signature: "+err.Error(), http.StatusBadRequest)
			return
		}

		if !CheckMAC(data, sigBytes, secretKey) {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, "Invalid X-Hub-Signature", http.StatusBadRequest)
		return
	}

	name := payload.Repository.Description
	url := payload.Repository.URL

	fmt.Println(event, sig, name, url)
}

func main() {
	secretKey = []byte(os.Getenv("GITHUBSECRET"))
	if len(secretKey) == 0 {
		panic("GITHUBSECRET environment variable not set")
	}

	http.HandleFunc("/githubnotify/", githubnotify)
	http.ListenAndServe(":8080", nil)
}
