package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Login - Exibe a página de login e processa o login de usuários
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Renderiza a página de login quando o método é GET
			tmpl := template.Must(template.ParseFiles("web/templates/login.html"))
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			// Processa o formulário de login quando o método é POST
			name := r.FormValue("name")
			password := r.FormValue("password")

			var storedPassword string
			err := db.QueryRow("SELECT password FROM users WHERE name = ?", name).Scan(&storedPassword)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			// Comparar a senha usando bcrypt
			err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
			if err != nil {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}

			// Gerar token JWT
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"name": name,
				"exp":  time.Now().Add(time.Hour * 1).Unix(),
			})

			tokenString, err := token.SignedString([]byte("lpkgSpxZw3jf2gmri/obJUry5QW7NZlC4QStyc0Cd/E="))
			if err != nil {
				http.Error(w, "Error generating token", http.StatusInternalServerError)
				return
			}

			// Define o cookie com o token JWT
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: time.Now().Add(time.Hour * 1),
				Path:    "/",
			})

			// Redireciona para a página tools
			http.Redirect(w, r, "/tools", http.StatusSeeOther)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

// Signup - Exibe a página de signup e processa o registro de novos usuários
func Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Renderiza a página de signup quando o método for GET
			tmpl := template.Must(template.ParseFiles("web/templates/signup.html"))
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			// Pega o nome de usuário e senha do formulário
			name := r.FormValue("name")
			email := r.FormValue("email")
			password := r.FormValue("password")

			// Verifica se o nome de usuário já existe
			var exists int
			err := db.QueryRow("SELECT COUNT(*) FROM users WHERE name = ?", name).Scan(&exists)
			if err != nil {
				http.Error(w, "Error checking name", http.StatusInternalServerError)
				return
			}

			if exists > 0 {
				http.Error(w, "Username already taken", http.StatusConflict)
				return
			}

			// Gera o hash da senha
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Error encrypting password", http.StatusInternalServerError)
				return
			}

			// Insere o novo usuário no banco de dados
			query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
			_, err = db.Exec(query, name, email, hashedPassword)
			if err != nil {
				http.Error(w, "Error creating user", http.StatusInternalServerError)
				return
			}

			// Redireciona para a página de login após o registro bem-sucedido
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

// Logout - Remove o token JWT (fazendo logout)
func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Definir o cookie de token com valor vazio e data de expiração passada
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Now().Add(-1 * time.Hour), // Define o cookie como expirado
			Path:    "/",
		})

		// Redireciona para a página de login após o logout
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
