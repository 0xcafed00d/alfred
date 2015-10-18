package main

import (
	//"fmt"
	"encoding/json"
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

type BuildInfo struct {
	PkgName         string
	GitHash         string
	BuildOK         bool
	TestOK          bool
	CoverageOK      bool
	CoveragePercent int
}

func buildPackage(pkg string) (binfo BuildInfo, err error) {
	binfo.PkgName = pkg

	log.Println("Processing Package:", pkg)
	err = goget(pkg, "build.log", &binfo)

	log.Println(" Running Tests on:", pkg)
	err = gotest(pkg, "test.log", &binfo)

	log.Println(" Processing Coverage on:", pkg)
	err = gocover(pkg, "cover.log", &binfo)

	return
}

func builder(queue <-chan string) {
	for {
		pkg := <-queue

		binfo, err := buildPackage(pkg)

		if err == nil {
			_, f, err := makePaths(pkg, "status.json")
			defer f.Close()

			enc := json.NewEncoder(f)
			err = enc.Encode(&binfo)
			if err != nil {
				log.Println(" Failed to write json status file", err)
			}

			log.Println(" Done", binfo)
		} else {
			log.Println(" Error: ", err, binfo)
		}
	}
}
