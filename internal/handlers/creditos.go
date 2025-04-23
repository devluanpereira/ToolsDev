package handlers

import (
	"consulta-cep/internal/utils"
	"database/sql"
	"encoding/json"
	"net/http"
)

// handlers/creditos.go
func VerificarCreditos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.GetUserIDFromRequest(r) // função que extrai o userID do JWT
		if err != nil {
			http.Error(w, "Não autenticado", http.StatusUnauthorized)
			return
		}

		var credits int
		err = db.QueryRow("SELECT credits FROM users WHERE id = ?", userID).Scan(&credits)
		if err != nil {
			http.Error(w, "Erro ao buscar créditos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"credits": credits,
		})
	}
}
