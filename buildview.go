package main

import (
	"html/template"
	"log"
	"net/http"
)

const buildsTemplateSrc = `
<html>
	<head>
		<title>Alfred Builds</title>
	</head>
	<body>
		<table>
			{{range .}}
  			<tr>
  				<td>{{.PkgHash}}</td>
    			<td>{{.PkgName}}</td>
    			<td>{{.BuildOK}}</td>
    			<td>{{.TestOK}}</td>
    			<td>{{.CoverageOK}}</td>
    			<td>{{.CoveragePercent}}</td>
  			</tr>
  			{{end}}
 		</table>
	</body>
</html>
`

var (
	buildsTemplate *template.Template
)

func init() {
	buildsTemplate = template.Must(template.New("builds").Parse(buildsTemplateSrc))

}

func buildView(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		return
	}

	binfos := doLs()

	log.Println(buildsTemplate.Execute(w, binfos))

}
