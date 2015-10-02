package main

import (
	//"fmt"
	"log"
)

type BuildQueue struct {
	queue chan string
}

func MakeBuildQueue() *BuildQueue {
	bq := BuildQueue{}
	bq.queue = make(chan string, 100)
	go builder(bq.queue)
	return &bq
}

func (bq *BuildQueue) EnQueue(pkg string) {
	bq.queue <- pkg
}

func builder(queue <-chan string) {
	for {
		pkg := <-queue

		log.Println("Processing Package: ", pkg)

		var err error
		err = goget(pkg, "build.log")
		if err == nil {
			log.Println(" Running Tests on: ", pkg)
			err = gotest(pkg, "test.log")
		}
		if err == nil {
			log.Println(" Processing Coverage on: ", pkg)
			err = gotest(pkg, "cover.log")
		}

		if err == nil {
			log.Println(" Done")
		} else {
			log.Println(" Error: ", err)
		}
	}
}
