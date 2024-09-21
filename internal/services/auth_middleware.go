package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Protected - Middleware para proteger rotas
func Protected(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Busca o cookie com o token
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			// Se o cookie não existir ou estiver vazio, redireciona para login
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
			return []byte("lpkgSpxZw3jf2gmri/obJUry5QW7NZlC4QStyc0Cd/E="), nil
		})

		if err != nil {
			// Se o token for inválido, redireciona para login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp, ok := claims["exp"].(float64)
			if !ok || int64(exp) < time.Now().Unix() {
				// Se o token estiver expirado, redireciona para login
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		} else {
			// Se o token for inválido, redireciona para login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Se o token for válido, chama o próximo handler (a rota protegida)
		next(w, r)
	}
}
