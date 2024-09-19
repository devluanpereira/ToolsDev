package handlers

import (
	"html/template"
	"net/http"
)

func InitHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/tools.html"))
	tmpl.Execute(w, nil)
}
