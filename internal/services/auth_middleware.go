package services

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Protected - Middleware para proteger rotas
func Protected(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Busca o cookie com o token
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value

		// Verifica e valida o token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verifica o algoritmo de assinatura
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo inesperado: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp, ok := claims["exp"].(float64)
			if !ok || int64(exp) < time.Now().Unix() {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Se o token for válido, chama o próximo handler (a rota protegida)
		next(w, r)
	}
}
