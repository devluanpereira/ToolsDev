package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"consulta-cep/internal/services"
)

// =======================================================================================================//
// HomeHandler renderiza a página inicial.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Corrigido: atribuição correta usando :=
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "Erro ao carregar a página: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// ======================================================================================================//
// CepHandler lida com a requisição de busca de CEP.
func CepHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "CEP não informado", http.StatusBadRequest)
		return
	}

	data, err := services.FetchCepData(cep)
	if err != nil {
		http.Error(w, "Erro ao buscar CEP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

//======================================================================================================//
