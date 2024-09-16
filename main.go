package main

import (
	"consulta-cep/internal/handlers"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/fatih/color"
)

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
	// Servindo arquivos estáticos (CSS, JS, imagens)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/buscar-cep", handlers.CepHandler)
	http.HandleFunc("/buscar-cnpj", handlers.CnpjHandler)
	http.HandleFunc("/buscar-code", handlers.BankHandler)
	http.HandleFunc("/iplookup", handlers.IPHandler)

	// Obtendo o IP local da máquina
	localIP := getLocalIP()

	// Exibindo o IP real da máquina no terminal
	color.Green(fmt.Sprintf("Servidor rodando em http://%s:8000", localIP))

	// Iniciando o servidor na porta 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}
