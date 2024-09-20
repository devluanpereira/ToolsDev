package main

import (
	"consulta-cep/internal/handlers"
	"consulta-cep/internal/services"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/fatih/color"
)

var db *sql.DB

// Função para obter o IP da máquina local
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "Erro ao obter IP"
	}

	for _, addr := range addrs {
		// Checa se é uma interface de rede com IP válido
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String() // Retorna o primeiro IP encontrado
		}
	}
	return "IP não encontrado"
}

func main() {

	var err error

	db, err = sql.Open("mysql", "projetosdev:luan@tcp(127.0.0.1:3306)/testeprimario")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Servindo arquivos estáticos (CSS, JS, imagens)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/signup", handlers.Signup(db))
	http.HandleFunc("/login", handlers.Login(db))
	http.HandleFunc("/logout", handlers.Logout())
	http.HandleFunc("/buscar-cep", handlers.CepHandler)
	http.HandleFunc("/buscar-cnpj", handlers.CnpjHandler)
	http.HandleFunc("/buscar-code", handlers.BankHandler)
	// Rotas protegidas pelo middleware
	http.HandleFunc("/tools", services.Protected(handlers.Tools))
	http.HandleFunc("/iplookup", services.Protected(handlers.IpLookup))

	// Obtendo o IP local da máquina
	localIP := getLocalIP()

	// Exibindo o IP real da máquina no terminal
	color.Green(fmt.Sprintf("Servidor rodando em http://%s:8000", localIP))

	// Iniciando o servidor na porta 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}
