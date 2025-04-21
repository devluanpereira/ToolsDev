package services

import (
	"consulta-cep/internal/models"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrMissingJWTSecret = errors.New("JWT_SECRET environment variable not set")

// HashPassword - Gera um hash da senha
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// CheckPasswordHash - Verifica se a senha está correta
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT - Gera o token JWT para o usuário
func GenerateJWT(user *models.User) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		return "", ErrMissingJWTSecret
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}
