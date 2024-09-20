package handlers

import (
	"consulta-cep/internal/services"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func GetIPFromRequest(r *http.Request) string {
	fowarded := r.Header.Get("X-Forwarded-For")
	if fowarded != "" {
		return strings.Split(fowarded, ",")[0]
	}

	return r.RemoteAddr
}

func IPHandler(w http.ResponseWriter, r *http.Request) {
	ip := "8.8.8.8" //Descomente se estiver em modo desenvolvimento
	//ip := GetIPFromRequest(r)

	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	ipInfo, err := services.GetIPInfo(ip)
	if err != nil {
		fmt.Println("Erro ao consultar IP:", err)
		http.Error(w, "Erro ao consultar IP", http.StatusInternalServerError)
		return
	}

	tmpl, _ := template.ParseFiles("web/templates/ip_lookup.html")
	tmpl.Execute(w, ipInfo)
}
