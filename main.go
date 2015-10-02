package main

import (
	"net/http"
	"os"
)

type BuildQueuer interface {
	EnQueue(pkg string)
}

func main() {
	secretKey := []byte(os.Getenv("GITHUBSECRET"))
	if len(secretKey) == 0 {
		panic("GITHUBSECRET environment variable not set")
	}

	webhook := GitHubWebHook{
		secretKey: secretKey,
		queuer:    MakeBuildQueue(),
	}

	http.HandleFunc("/githubnotify/", webhook.notify)
	http.ListenAndServe(":8080", nil)
}
