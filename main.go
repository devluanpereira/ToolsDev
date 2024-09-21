package main

import (
	"consulta-cep/internal/handlers"
	"consulta-cep/internal/services"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Carregando as variáveis do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	// Inicializando a conexão com o banco de dados
	db, err := initDB()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()
	log.Println("Conexão com o banco de dados estabelecida com sucesso.")

	// Servindo arquivos estáticos (CSS, JS, imagens)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Definindo rotas
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/signup", handlers.Signup(db))
	http.HandleFunc("/login", handlers.Login(db))
	http.HandleFunc("/logout", handlers.Logout())
	http.HandleFunc("/buscar-cep", handlers.CepHandler)
	http.HandleFunc("/buscar-cnpj", handlers.CnpjHandler)
	http.HandleFunc("/buscar-code", handlers.BankHandler)

	// Rotas protegidas pelo middleware
	http.HandleFunc("/tools", services.Protected(handlers.Tools))

	// Inicializando o servidor na porta 8000
	server := &http.Server{
		Addr:         ":8000",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Println("Servidor rodando na porta 8000.")
	log.Fatal(server.ListenAndServe())
}

// initDB inicializa a conexão com o banco de dados MySQL
func initDB() (*sql.DB, error) {
	// Obtendo as variáveis de ambiente
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Configurando a string de conexão
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Inicializando a conexão com o banco de dados
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir a conexão com o banco de dados: %w", err)
	}

	// Verificando a conexão com o banco de dados
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	return db, nil
}
