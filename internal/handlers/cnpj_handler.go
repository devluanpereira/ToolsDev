package handlers

import (
	"consulta-cep/internal/services"
	"encoding/json"
	"net/http"
)

func CnpjHandler(w http.ResponseWriter, r *http.Request) {
	cnpj := r.URL.Query().Get("cnpj")
	if cnpj == "" {
		http.Error(w, "CNPJ não informado", http.StatusBadRequest)
		return
	}

	data, err := services.FetchCnpjData(cnpj)
	if err != nil {
		http.Error(w, "Erro ao buscar CNPJ: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
