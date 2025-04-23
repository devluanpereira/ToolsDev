package utils

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("sua-chave-secreta")

func GetUserIDFromRequest(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("token ausente")
	}

	tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if idFloat, ok := claims["id"].(float64); ok {
			return int(idFloat), nil
		}
	}
	return 0, err
}
