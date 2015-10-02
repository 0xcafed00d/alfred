package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func generatePackageHash(pkg string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(pkg)))
}

func makePaths(pkg, logfile string) (gopath string, logwriter io.WriteCloser, err error) {

	gopath, err = filepath.Abs(generatePackageHash(pkg))
	if err != nil {
		return
	}

	logwriter, err = os.Create(filepath.Join(gopath, logfile))
	if err != nil {
		return
	}

	return
}

func goget(pkg, logfile string) error {

	err := os.MkdirAll(generatePackageHash(pkg), os.ModePerm)
	if err != nil {
		return err
	}

	gopath, logwriter, err := makePaths(pkg, logfile)
	if err != nil {
		return err
	}
	defer logwriter.Close()

	return execWithTimeout("go", "get -v -u -t "+pkg+"/...", gopath, logwriter, 300*time.Second)
}

func gotest(pkg, logfile string) error {

	gopath, logwriter, err := makePaths(pkg, logfile)
	if err != nil {
		return err
	}
	defer logwriter.Close()

	coverdata := filepath.Join(gopath, "coverdata.out")

	args := fmt.Sprintf("test -v -covermode=count -coverprofile=%s %s", coverdata, pkg)
	return execWithTimeout("go", args, gopath, logwriter, 300*time.Second)
}

func gocover(pkg, logfile string) error {

	gopath, logwriter, err := makePaths(pkg, logfile)
	if err != nil {
		return err
	}
	defer logwriter.Close()

	coverdata := filepath.Join(gopath, "coverdata.out")
	html := filepath.Join(gopath, "coverdata.html")

	args := fmt.Sprintf("tool cover -html=%s -o %s", coverdata, html)
	return execWithTimeout("go", args, gopath, logwriter, 300*time.Second)
}
