package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MercadoPagoResponse struct {
	InitPoint string `json:"init_point"`
}

func extractUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("token ausente")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return 0, fmt.Errorf("formato do token inválido")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("token inválido")
	}

	idFloat, ok := claims["userID"].(float64)
	if !ok {
		return 0, fmt.Errorf("ID do usuário não encontrado no token")
	}

	return int(idFloat), nil
}

func CriarPagamento(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Renderiza o formulário de pagamento
			tmpl := template.Must(template.ParseFiles("web/templates/criar_pagamento.html"))
			tmpl.Execute(w, nil)

		case http.MethodPost:
			// Processa o pagamento
			_, err := extractUserIDFromToken(r)
			if err != nil {
				http.Error(w, "Usuário não autenticado: "+err.Error(), http.StatusUnauthorized)
				return
			}

			if err := r.ParseForm(); err != nil {
				http.Error(w, "Erro ao analisar formulário", http.StatusBadRequest)
				return
			}

			email := r.FormValue("email")
			quantidadeStr := r.FormValue("quantidade")

			quantidade, err := strconv.Atoi(quantidadeStr)
			if err != nil || quantidade <= 0 {
				http.Error(w, "Quantidade inválida", http.StatusBadRequest)
				return
			}

			// Criação do pagamento via Mercado Pago
			preference := map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"title":       "Recarga de Créditos",
						"description": fmt.Sprintf("Recarga de %d créditos", quantidade),
						"quantity":    1,
						"unit_price":  float64(quantidade),
						"currency_id": "BRL",
					},
				},
				"payer": map[string]string{
					"email": email,
				},
				"back_urls": map[string]string{
					"success": "http://localhost:8000/pagamento-sucesso",
					"failure": "http://localhost:8000/pagamento-falhou",
				},
				"auto_return": "approved",
			}

			payload, _ := json.Marshal(preference)

			req, err := http.NewRequest("POST", "https://api.mercadopago.com/checkout/preferences", bytes.NewBuffer(payload))
			if err != nil {
				http.Error(w, "Erro ao criar requisição", http.StatusInternalServerError)
				return
			}
			req.Header.Set("Authorization", "Bearer "+os.Getenv("MP_ACCESS_TOKEN"))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, "Erro ao enviar requisição ao MercadoPago", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			var mpResponse MercadoPagoResponse
			if err := json.NewDecoder(resp.Body).Decode(&mpResponse); err != nil {
				http.Error(w, "Erro ao ler resposta do MercadoPago", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, mpResponse.InitPoint, http.StatusSeeOther)

		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}
}
