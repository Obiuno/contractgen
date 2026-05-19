package web

import (
	"embed"
	"html/template"
)

//go:embed templates/*.html
var templateFS embed.FS

var (
	IndexTmpl   = template.Must(template.ParseFS(templateFS, "templates/index.html"))
	ResultsTmpl = template.Must(template.ParseFS(templateFS, "templates/results.html"))
)