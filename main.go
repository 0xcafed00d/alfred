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

func getSecretKey() ([]byte, error) {
	secretKey := []byte(os.Getenv("GITHUBSECRET"))
	if len(secretKey) == 0 {
		return nil, errors.New("GITHUBSECRET environment variable not set")
	}
	return secretKey, nil
}

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
	http.ListenAndServe(":8080", nil)
}

func cmdKick(args []string) {

}

func runCommand(cmd string, args []string) {
	switch cmd {
	case "serve":
		cmdServe(args)
	case "kick":
		cmdKick(args)
	}
}

func main() {
	if len(os.Args) > 1 {
		runCommand(os.Args[1], os.Args[2:])
	}
}
