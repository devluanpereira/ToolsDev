package handlers

import (
	"html/template"
	"net/http"
)

func ExibirFormularioPagamento() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("web/templates/criar_pagamento.html"))

		data := struct {
			UserID int
			Email  string
		}{
			UserID: 1, // Substituir com dados reais
			Email:  "teste@exemplo.com",
		}

		tmpl.Execute(w, data)
	}
}
