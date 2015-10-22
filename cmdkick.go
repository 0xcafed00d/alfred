package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func cmdKick(args []string) {
	if len(args) != 2 {
		exitOnError(errors.New("Usage: alfred kick <go pkgname> <serverURL>"))
	}

	pkg := args[0]
	url := args[1]

	secretKey, err := getSecretKey()
	exitOnError(err)

	payload := GitHubPayload{}

	payload.Repository.URL = "https://" + pkg

	plbytes, err := json.Marshal(&payload)
	exitOnError(err)

	req, err := http.NewRequest("POST", url+"/githubnotify/",
		bytes.NewBuffer(plbytes))
	exitOnError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature", "sha1="+generateHMAC(plbytes, secretKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	exitOnError(err)
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
