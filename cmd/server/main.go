package main

import (
	"consulta-cep/internal/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/buscar-cep", handlers.CepHandler)
	http.HandleFunc("/buscar-cnpj", handlers.CnpjHandler)

	fmt.Println("Servidor rodando em http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
