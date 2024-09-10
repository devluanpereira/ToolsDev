package main

import (
	"consulta-cep/internal/handlers"
	"fmt"
	"log"
	"net/http"
)

// =========================================================================//
func main() {

	// Servindo arquivos estáticos (CSS, JS, imagens)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/buscar-cep", handlers.CepHandler)
	http.HandleFunc("/buscar-cnpj", handlers.CnpjHandler)
	http.HandleFunc("/buscar-code", handlers.BankHandler)

	fmt.Println("Servidor rodando em http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

//=========================================================================//
