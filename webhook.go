package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type GitHubWebHook struct {
	secretKey []byte
	queuer    BuildQueuer
}

func generateHMAC(data, key []byte) string {
	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

func (wh *GitHubWebHook) verifyHMAC(w http.ResponseWriter, r *http.Request, body []byte) bool {
	sig := r.Header.Get("X-Hub-Signature")

	expectedMAC := "sha1=" + generateHMAC(body, wh.secretKey)
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
		r.Body.Close()

		if !wh.verifyHMAC(w, r, body) {
			return
		}

		payload := GitHubPayload{}
		err = json.Unmarshal(body, &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pkg := strings.TrimPrefix(payload.Repository.URL, "https://")

		wh.queuer.EnQueue(pkg)
	}
}
