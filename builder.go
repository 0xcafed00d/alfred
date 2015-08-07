package main

import (
	"fmt"
)

func builder(repos <-chan string) {
	for {
		repo := <-repos
		fmt.Println(repo)
	}
}
