package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func boolOk(ok bool) string {
	if ok {
		return "OK    "
	} else {
		return "Failed"
	}
}

func printBuildInfo(location string, binfo *BuildInfo) {
	fmt.Printf("%s %s\n", location, binfo.PkgName)
	fmt.Printf("Build: %s  Test: %s  Coverage: %s %d%%\n\n",
		boolOk(binfo.BuildOK),
		boolOk(binfo.TestOK),
		boolOk(binfo.CoverageOK),
		binfo.CoveragePercent)
}

func doLs() (binfos []BuildInfo) {
	matches, err := filepath.Glob("*")
	exitOnError(err)

	for _, v := range matches {
		path := filepath.Join(v, "status.json")

		f, err := os.Open(path)
		if err != nil {
			continue
		}

		status, err := ioutil.ReadAll(f)
		if err != nil {
			continue
		}

		binfo := BuildInfo{}
		err = json.Unmarshal(status, &binfo)
		if err != nil {
			continue
		}
		binfos = append(binfos, binfo)
	}
	return
}

func cmdLs(args []string) {
	err := os.Chdir("data")
	exitOnError(err)

	matches, err := filepath.Glob("*")
	exitOnError(err)

	for _, v := range matches {
		path := filepath.Join(v, "status.json")

		f, err := os.Open(path)
		if err != nil {
			continue
		}

		status, err := ioutil.ReadAll(f)
		if err != nil {
			continue
		}

		binfo := BuildInfo{}
		err = json.Unmarshal(status, &binfo)
		if err != nil {
			continue
		}

		printBuildInfo(v, &binfo)
	}
}
