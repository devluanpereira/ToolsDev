package handlers

import (
	"consulta-cep/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Estrutura para passar dados ao template
type LoginData struct {
	ErrorMessage string
}

// Login - Exibe a página de login e processa o login de usuários
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handler Login chamado!") // Log no início

		if r.Method == http.MethodGet {
			// Renderiza a página de login quando o método é GET
			tmpl := template.Must(template.ParseFiles("web/templates/login.html"))
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			fmt.Println("Método POST recebido!") // Log para método POST

			email := r.FormValue("email")
			password := r.FormValue("password")
			fmt.Println("Email:", email, "Senha:", password) // Log dos dados do formulário

			var storedPassword string
			var user models.User
			err := db.QueryRow("SELECT id, password, role, credits FROM users WHERE email = ?", email).Scan(&user.ID, &storedPassword, &user.Role, &user.Credits)
			if err != nil {
				fmt.Println("Erro na consulta:", err) // Log do erro na consulta
				tmpl := template.Must(template.ParseFiles("web/templates/login.html"))
				tmpl.Execute(w, LoginData{ErrorMessage: "Email ou senha inválidos."})
				return
			}

			// Comparar a senha usando bcrypt
			err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
			if err != nil {
				fmt.Println("Erro na comparação da senha:", err) // Log do erro na comparação
				tmpl := template.Must(template.ParseFiles("web/templates/login.html"))
				tmpl.Execute(w, LoginData{ErrorMessage: "Email ou senha inválidos."})
				return
			}

			// Gerar token JWT
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"email":   email,
				"userID":  user.ID,
				"role":    user.Role,
				"credits": user.Credits,
				"exp":     time.Now().Add(time.Hour * 1).Unix(),
			})

			tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))) // Usa a variável do .env
			if err != nil {
				fmt.Println("Erro ao gerar o token JWT:", err) // Log do erro na geração do token
				http.Error(w, "Error generating token", http.StatusInternalServerError)
				return
			}

			// Define o cookie com o token JWT
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    tokenString,
				Path:     "/",
				HttpOnly: true,         // Adiciona HttpOnly para segurança
				Secure:   r.TLS != nil, // Define Secure apenas se a conexão for HTTPS
			})

			fmt.Println("Login bem-sucedido, redirecionando para /tools") // Log de sucesso

			// Redireciona para a página tools
			http.Redirect(w, r, "/tools", http.StatusSeeOther)
		} else {
			fmt.Println("Método inválido:", r.Method) // Log para método inválido
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
			// Pega o nome de usuário, email e senha do formulário
			name := r.FormValue("name")
			email := r.FormValue("email")
			password := r.FormValue("password")

			// Verifica se o email já existe
			var exists int
			err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&exists)
			if err != nil {
				fmt.Println("Erro ao verificar email:", err)
				http.Error(w, "Error checking email", http.StatusInternalServerError)
				return
			}

			if exists > 0 {
				fmt.Println("Email já está em uso:", email)
				tmpl := template.Must(template.ParseFiles("web/templates/signup.html"))
				tmpl.Execute(w, struct {
					ErrorMessage string
				}{ErrorMessage: "Email já está em uso."})
				return
			}

			// Gera o hash da senha
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Println("Erro ao criptografar senha:", err)
				http.Error(w, "Error encrypting password", http.StatusInternalServerError)
				return
			}

			// Insere o novo usuário no banco de dados
			query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
			_, err = db.Exec(query, name, email, hashedPassword)
			if err != nil {
				fmt.Println("Erro ao criar usuário:", err)
				http.Error(w, "Error creating user", http.StatusInternalServerError)
				return
			}

			fmt.Println("Usuário criado com sucesso, redirecionando para /login")
			// Redireciona para a página de login após o registro bem-sucedido
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			fmt.Println("Método inválido no signup:", r.Method)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

// Logout - Remove o token JWT (fazendo logout)
func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Definir o cookie de token com valor vazio e data de expiração passada
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour), // Define o cookie como expirado
			Path:     "/",
			HttpOnly: true,
			Secure:   r.TLS != nil,
		})

		fmt.Println("Logout realizado, redirecionando para /login")
		// Redireciona para a página de login após o logout
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// Funcao para pegar o User ID do JWT
func GetUserIDFromRequest(r *http.Request) (int, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("Token não encontrado na requisição:", err)
		return 0, errors.New("token não encontrado")
	}

	tokenString := cookie.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		fmt.Println("Token JWT inválido:", err)
		return 0, errors.New("token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Claims JWT inválidas")
		return 0, errors.New("claims invalidas")
	}

	uid, ok := claims["userID"].(float64)
	if !ok {
		fmt.Println("userID ausente ou com tipo incorreto no token")
		return 0, errors.New("userID ausente no token")
	}

	fmt.Println("UserID extraído do token:", int(uid))
	return int(uid), nil
}

func GetRoleFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("Token não encontrado na requisição (GetRole):", err)
		return "", errors.New("token não encontrado")
	}

	tokenString := cookie.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		fmt.Println("Token JWT inválido (GetRole):", err)
		return "", errors.New("token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Claims JWT inválidas (GetRole)")
		return "", errors.New("claims inválidas")
	}

	role, ok := claims["role"].(string)
	if !ok {
		fmt.Println("Role ausente ou com tipo incorreto no token")
		return "", errors.New("role ausente no token")
	}

	fmt.Println("Role extraída do token:", role)
	return role, nil
}

func ConsumeCredit(db *sql.DB, userID int) error {
	var credits int
	err := db.QueryRow("SELECT credits FROM users WHERE id = ?", userID).Scan(&credits)
	if err != nil {
		fmt.Println("Erro ao obter créditos do usuário:", err)
		return err
	}

	if credits <= 0 {
		fmt.Println("Sem créditos disponíveis para o usuário:", userID)
		return errors.New("Sem créditos disponiveis")
	}

	_, err = db.Exec("UPDATE users SET credits = credits - 1 WHERE id = ?", userID) // Corrigido para subtrair 1
	if err != nil {
		fmt.Println("Erro ao consumir crédito do usuário:", err)
		return err
	}

	fmt.Println("Crédito consumido com sucesso para o usuário:", userID)
	return nil
}

func ConsumeCreditCep(db *sql.DB, userID int) error {
	var credits int
	err := db.QueryRow("SELECT credits FROM users WHERE id = ?", userID).Scan(&credits)
	if err != nil {
		fmt.Println("Erro ao obter créditos do usuário (ConsumeCreditCep):", err)
		return err
	}

	if credits <= 0 {
		fmt.Println("Sem créditos disponíveis para o usuário (ConsumeCreditCep):", userID)
		return errors.New("Sem créditos disponiveis")
	}

	_, err = db.Exec("UPDATE users SET credits = credits - 1 WHERE id = ?", userID) // Assumindo que a consulta de CEP também consome 1 crédito
	if err != nil {
		fmt.Println("Erro ao consumir crédito do usuário (ConsumeCreditCep):", err)
		return err
	}

	fmt.Println("Crédito consumido para consulta de CEP do usuário:", userID)
	return nil
}
