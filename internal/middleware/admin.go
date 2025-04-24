package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func AdminOnlyMiddleware(next http.HandlerFunc, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, err := GetRoleFromRequest(r)
		if err != nil || role != "admin" {
			http.Error(w, "Acesso negado", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func GetRoleFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", errors.New("token não encontrado")
	}

	tokenString := cookie.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("claims inválidas")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("role ausente no token")
	}

	return role, nil
}
