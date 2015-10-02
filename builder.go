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
		var err error

		run := func(f func(pkg, logfile string) error, logfile, logmsg string) {
			if err == nil {
				err = f(pkg, logfile)
				log.Println(logmsg, pkg)
			}
		}

		run(goget, "build.log", "Processing Package:")
		run(gotest, "test.log", " Running Tests on:")
		run(gocover, "cover.log", " Processing Coverage on:")

		if err == nil {
			log.Println(" Done")
		} else {
			log.Println(" Error: ", err)
		}
	}
}
