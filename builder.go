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
	PkgHash         string
	GitHash         string
	BuildOK         bool
	TestOK          bool
	CoverageOK      bool
	CoveragePercent int
}

func buildPackage(pkg string) (binfo BuildInfo) {
	binfo.PkgName = pkg
	binfo.PkgHash = generatePackageHash(pkg)

	log.Println("Processing Package:", pkg)
	err := goget(pkg, "build.log", &binfo)
	if err != nil {
		return
	}

	log.Println(" Running Tests on:", pkg)
	err = gotest(pkg, "test.log", &binfo)

	if err == nil {
		log.Println(" Processing Coverage on:", pkg)
		gocover(pkg, "cover.log", &binfo)
	}

	return
}

func builder(queue <-chan string) {
	for {
		pkg := <-queue

		binfo := buildPackage(pkg)

		_, f, err := makePaths(pkg, "status.json")
		defer f.Close()

		enc := json.NewEncoder(f)
		err = enc.Encode(&binfo)
		if err != nil {
			log.Println(" Failed to write json status file", err)
		}

		log.Println(" Done", binfo)
	}
}
