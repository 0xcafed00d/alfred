package main

import (
	"fmt"
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
