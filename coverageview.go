package main

import (
	"net/http"
	"path/filepath"
)

func coverageView(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	pkgHash := r.URL.Query().Get("pkg")

	http.ServeFile(w, r, filepath.Join(pkgHash, "coverdata.html"))
}
