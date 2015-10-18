package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

type BuildQueuer interface {
	EnQueue(pkg string)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runCommand(cmd string, args []string) {
	fmt.Println(cmd, args)
}

func main() {
	secretKey := []byte(os.Getenv("GITHUBSECRET"))
	if len(secretKey) == 0 {
		exitOnError(errors.New("GITHUBSECRET environment variable not set"))
	}

	if len(os.Args) > 1 {
		runCommand(os.Args[1], os.Args[2:])
		return
	}

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
	http.ListenAndServe(":8080", nil)
}
