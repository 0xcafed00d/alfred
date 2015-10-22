package main

import (
	"errors"
	"net/http"
	"os"
)

func cmdServe(args []string) {
	secretKey, err := getSecretKey()
	exitOnError(err)

	finfo, err := os.Stat("data")
	if err == nil {
		if !finfo.IsDir() {
			exitOnError(errors.New("'data' present but not directory"))
		}
	} else {
		err = os.Mkdir("data", os.ModePerm)
		exitOnError(err)
	}

	err = os.Chdir("data")
	exitOnError(err)

	webhook := GitHubWebHook{
		secretKey: secretKey,
		queuer:    MakeBuildQueue(),
	}

	http.HandleFunc("/githubnotify/", webhook.notify)
	exitOnError(http.ListenAndServe(":8080", nil))
}
