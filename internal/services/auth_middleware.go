package services

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Protected - Middleware para proteger rotas
func Protected(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Busca o cookie com o token
		cookie, err := r.Cookie("token")
		if err != nil {
			// Se o cookie não existir, redireciona para a página de login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value

		// Verifica e valida o token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("lpkgSpxZw3jf2gmri/obJUry5QW7NZlC4QStyc0Cd/E="), nil
		})

		if err != nil {
			// Se o token for inválido, redireciona para login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := int64(claims["exp"].(float64))
			if exp < time.Now().Unix() {
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
