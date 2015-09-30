package main

import (
	"crypto/sha1"
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

func generatePackageHash(pkg string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(pkg)))
}

func goget(pkg, logfile string) error {

	gopath := generatePackageHash(pkg)

	err := os.MkdirAll(gopath, os.ModePerm)
	if err != nil {
		return err
	}

	gopath, err = filepath.Abs(gopath)
	if err != nil {
		return err
	}

	log, err := os.Create(filepath.Join(gopath, logfile))
	if err != nil {
		return err
	}
	defer log.Close()

	err = execWithTimeout("go", "get -v -u -t "+pkg+"/...", gopath, log, 300*time.Second)
	return err
}

func gotest(pkg, logfile string) error {

	gopath := generatePackageHash(pkg)
	gopath, err := filepath.Abs(gopath)
	if err != nil {
		return err
	}

	log, err := os.Create(filepath.Join(gopath, logfile))
	if err != nil {
		return err
	}
	defer log.Close()

	coverdata := filepath.Join(gopath, "coverdata.out")

	args := fmt.Sprintf("test -v -covermode=count -coverprofile=%s %s", coverdata, pkg)
	err = execWithTimeout("go", args, gopath, log, 300*time.Second)
	return err
}

func gocover(pkg, logfile string) error {

	gopath := generatePackageHash(pkg)
	gopath, err := filepath.Abs(gopath)
	if err != nil {
		return err
	}

	log, err := os.Create(filepath.Join(gopath, logfile))
	if err != nil {
		return err
	}
	defer log.Close()

	coverdata := filepath.Join(gopath, "coverdata.out")
	html := filepath.Join(gopath, "coverdata.html")

	args := fmt.Sprintf("tool cover -html=%s -o %s", coverdata, html)
	err = execWithTimeout("go", args, gopath, log, 300*time.Second)
	return err
}
