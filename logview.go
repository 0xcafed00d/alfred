package main

import (
	"html/template"
	"log"
	"net/http"
)

const logTemplateSrc = `
<html>
	<head>
		<title>Alfred {{.PkgName}} Log</title>
	</head>
	<body>
	</body>
</html>
`

var (
	logTemplate *template.Template
)

func init() {
	logTemplate = template.Must(template.New("log").Parse(logTemplateSrc))

}

func logView(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	log.Println(r.URL.Query())
}
