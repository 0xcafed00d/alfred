package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type LogInfo struct {
	PkgHash string
	LogType string
	LogBody string
}

const logTemplateSrc = `
<html>
	<head>
		<title>Alfred {{.PkgHash}} {{.LogType}} Log</title>
		Alfred {{.PkgHash}} {{.LogType}} log
	</head>
	<body>
		<pre>
			{{.LogBody}}
		</pre>
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

	linfo := LogInfo{}

	linfo.PkgHash = r.URL.Query().Get("pkg")
	linfo.LogType = r.URL.Query().Get("type")

	body, err := ioutil.ReadFile(filepath.Join(linfo.PkgHash, linfo.LogType+".log"))
	linfo.LogBody = string(body)

	log.Println(err, logTemplate.Execute(w, linfo))
}
