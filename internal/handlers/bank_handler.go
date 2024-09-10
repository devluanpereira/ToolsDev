package handlers

//==============================================================================================//
import (
	"consulta-cep/internal/services"
	"encoding/json"
	"net/http"
)

//==============================================================================================//

func BankHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "ISPB não informado", http.StatusBadRequest)
		return
	}

	data, err := services.FetchBankData(code)
	if err != nil {
		http.Error(w, "Erro ao buscar code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

//=============================================================================================//
