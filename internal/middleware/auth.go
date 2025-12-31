package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
// AQUI COLOQUE O JWT_SECRET NO ENV
		tokenStr := cookie.Value
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID := int(userIDFloat)
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func VerificarCreditosMiddleware(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Aqui você pode pegar o ID da sessão ou do cookie
		userID := r.Context().Value("userID").(int)

		var credits int
		err := db.QueryRow("SELECT credits FROM users WHERE id = ?", userID).Scan(&credits)
		if err != nil {
			http.Error(w, "Erro ao verificar créditos", http.StatusInternalServerError)
			return
		}

		if credits <= 0 {
			// Redireciona para a tela de recarga
			http.Redirect(w, r, "/criar-pagamento", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
}
