package main

import (
	"net/http"
	"os"
)

func main() {
	secretKey := []byte(os.Getenv("GITHUBSECRET"))
	if len(secretKey) == 0 {
		panic("GITHUBSECRET environment variable not set")
	}

	builderChan := make(chan string)
	go builder(builderChan)

	webhook := GitHubWebHook{
		secretKey:   secretKey,
		builderChan: builderChan,
	}

	http.HandleFunc("/githubnotify/", webhook.notify)
	http.ListenAndServe(":8080", nil)
}
