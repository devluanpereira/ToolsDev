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
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established successfully.")

	// Create necessary tables
	if err := createTables(db); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	// Serve static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/signup", handlers.Signup(db))
	http.HandleFunc("/login", handlers.Login(db))
	http.HandleFunc("/logout", handlers.Logout())
	http.HandleFunc("/buscar-cep", handlers.CepHandler)
	http.HandleFunc("/buscar-cnpj", handlers.CnpjHandler)
	http.HandleFunc("/buscar-code", handlers.BankHandler)

	// Protected routes
	http.HandleFunc("/tools", services.Protected(handlers.Tools))

	server := &http.Server{
		Addr:         ":8000",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Println("Server running on port 8000.")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server falhou: %v", err)
	}
}

func initDB() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbPort == "" || dbName == "" {
		return nil, fmt.Errorf("one or more required environment variables are not set")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, os.Getenv("DB_HOST"), dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,  -- Removido UNIQUE
		email VARCHAR(255) NOT NULL UNIQUE,  -- Email deve ser único
		password VARCHAR(255) NOT NULL
	);`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}

	return nil
}
