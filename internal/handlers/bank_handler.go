package handlers

//==============================================================================================//
import (
	"consulta-cep/internal/services"
	"database/sql"
	"encoding/json"
	"net/http"
)

//==============================================================================================//

func BankHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "ISPB não informado", http.StatusBadRequest)
			return
		}

		userID, err := GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Usuario não autenticado: "+err.Error(), http.StatusUnauthorized)
			return
		}

		err = ConsumeCredit(db, userID)
		if err != nil {
			http.Error(w, "Error ao consumir crédito: "+err.Error(), http.StatusForbidden)
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
}

//=============================================================================================//
