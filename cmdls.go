package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func cmdLs(args []string) {
	err := os.Chdir("data")
	exitOnError(err)

	matches, err := filepath.Glob("*")
	exitOnError(err)

	for _, v := range matches {
		path := filepath.Join(v, "status.json")
		f, err := os.Open(path)
		if err == nil {
			status, err := ioutil.ReadAll(f)
			if err == nil {
				binfo := BuildInfo{}
				err = json.Unmarshal(status, &binfo)
				if err == nil {
					fmt.Println(binfo)
				}
			}
		}
	}
}
