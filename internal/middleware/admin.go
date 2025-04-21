package middleware

import (
	"database/sql"
	"net/http"
)

func AdminOnlyMiddleware(next http.HandlerFunc, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, err := GetRoleFromRequest(r, db)
		if err != nil || role != "admin" {
			http.Error(w, "Acesso negado", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func GetRoleFromRequest(r *http.Request, db *sql.DB) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}

	var role string
	err = db.QueryRow("SELECT role FROM users WHERE token = ?", cookie.Value).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}
