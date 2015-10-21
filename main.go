package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	payload := GitHubPayload{}

	payload.Repository.URL = "https://github.com/simulatedsimian/alfred"

	plbytes, err := json.Marshal(&payload)
	exitOnError(err)

	req, err := http.NewRequest("POST", "http://localhost:8080/githubnotify/",
		bytes.NewBuffer(plbytes))
	exitOnError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")

	client := &http.Client{}
	resp, err := client.Do(req)
	exitOnError(err)
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
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
