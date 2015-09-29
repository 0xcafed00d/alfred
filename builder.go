package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type KnownRepos struct {
	url string
}

func builder(repos <-chan string) {
	for {
		repo := <-repos
		fmt.Println(repo)
	}
}

func goget(repo, gopath, logfile string) error {
	err := os.MkdirAll(gopath, os.ModePerm)
	if err != nil {
		return err
	}

	gopath, err = filepath.Abs(gopath)
	if err != nil {
		return err
	}

	log, err := os.Create(filepath.Join(gopath, logfile))
	defer log.Close()

	if err != nil {
		return err
	}

	err = execWithTimeout("go", "get -v -u -t "+repo+"/...", gopath, log, 30*time.Second)
	return err
}
