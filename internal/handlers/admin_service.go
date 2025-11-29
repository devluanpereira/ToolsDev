package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
)

type User struct {
	ID      int
	Email   string
	Credits int
}

func AdicionarCredito(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("web/templates/admin/credito.html")
			if err != nil {
				http.Error(w, "Erro ao carregar template", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
			return 
		}

		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			quantidadeStr := r.FormValue("quantidade")

			quantidade, err := strconv.Atoi(quantidadeStr)
			if err != nil {
				http.Error(w, "Quantidade inválida", http.StatusBadRequest)
				return
			}

			_, err = db.Exec("UPDATE users SET credits = credits + ? WHERE email = ?", quantidade, email)
			if err != nil {
				http.Error(w, "Erro ao adicionar créditos", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/admin/adicionar", http.StatusSeeOther)
			return
		}

		// Se não for GET nem POST
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}
