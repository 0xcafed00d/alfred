package main

import (
	"errors"
	"fmt"
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
