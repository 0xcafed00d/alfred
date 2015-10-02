package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type GitHubWebHook struct {
	secretKey []byte
	queuer    BuildQueuer
}

func (wh *GitHubWebHook) verifyHMAC(w http.ResponseWriter, r *http.Request, body []byte) bool {
	sig := r.Header.Get("X-Hub-Signature")

	mac := hmac.New(sha1.New, wh.secretKey)
	mac.Write(body)
	expectedMAC := "sha1=" + hex.EncodeToString(mac.Sum(nil))
	equal := hmac.Equal([]byte(sig), []byte(expectedMAC))

	if !equal {
		http.Error(w, "X-Hub-Signature Mismatch", http.StatusUnauthorized)
	}

	return equal
}

func (wh *GitHubWebHook) notify(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Github-Event") == "push" {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !wh.verifyHMAC(w, r, body) {
			return
		}

		payload := GitHubPayload{}
		err = json.Unmarshal(body, &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		wh.queuer.EnQueue(payload.Repository.URL)
	}
}
