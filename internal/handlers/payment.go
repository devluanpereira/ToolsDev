package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MercadoPagoResponse struct {
	InitPoint string `json:"init_point"`
}

func extractUserIDFromCookie(r *http.Request) (int, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return 0, fmt.Errorf("cookie 'token' não encontrado")
	}

	tokenString := cookie.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("token JWT inválido no cookie: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("claims do token inválidas no cookie")
	}

	idFloat, ok := claims["userID"].(float64)
	if !ok {
		return 0, fmt.Errorf("ID do usuário não encontrado no token do cookie")
	}

	return int(idFloat), nil
}

func CriarPagamento(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tmpl := template.Must(template.ParseFiles("web/templates/criar_pagamento.html"))
			tmpl.Execute(w, nil)

		case http.MethodPost:
			userID, err := extractUserIDFromCookie(r)
			if err != nil {
				http.Error(w, "Usuário não autenticado: "+err.Error(), http.StatusUnauthorized)
				return
			}
			fmt.Printf("Usuário autenticado com ID (via cookie): %d\n", userID)

			if err := r.ParseForm(); err != nil {
				fmt.Printf("Erro ao analisar formulário: %v\n", err)
				http.Error(w, "Erro ao analisar formulário", http.StatusBadRequest)
				return
			}
			fmt.Printf("Dados do formulário após ParseForm: %+v\n", r.PostForm)

			email := r.FormValue("email")
			quantidadeStr := r.FormValue("quantidade")
			fmt.Printf("Valor da quantidade recebido: '%s'\n", quantidadeStr)

			quantidade, err := strconv.Atoi(quantidadeStr)
			if err != nil {
				fmt.Printf("Erro ao converter quantidade para inteiro: %v\n", err)
				http.Error(w, "Quantidade inválida", http.StatusBadRequest)
				return
			}
			if quantidade <= 0 {
				fmt.Println("Quantidade é zero ou negativa:", quantidade)
				http.Error(w, "Quantidade inválida", http.StatusBadRequest)
				return
			}

			fmt.Printf("Quantidade convertida para inteiro: %d\n", quantidade)

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
					"success": fmt.Sprintf("http://localhost:8000/pagamento-sucesso?user_id=%d&credits=%d", userID, quantidade),
					"failure": "http://localhost:8000/pagamento-falhou",
				},
				"metadata": map[string]interface{}{
					"user_id": userID,
					"credits": quantidade,
					"email":   email,
				},
			}

			payload, err := json.Marshal(preference)
			if err != nil {
				http.Error(w, "Erro ao criar payload JSON", http.StatusInternalServerError)
				return
			}
			fmt.Printf("Payload JSON para o Mercado Pago: %s\n", string(payload))

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

			if resp.StatusCode != http.StatusCreated {
				var bodyBytes []byte
				bodyBytes, _ = io.ReadAll(resp.Body)
				fmt.Printf("Erro do MercadoPago: Status Code %d, Body: %s\n", resp.StatusCode, string(bodyBytes))
				http.Error(w, "Erro na comunicação com o MercadoPago", http.StatusInternalServerError)
				return
			}

			var mpResponse MercadoPagoResponse
			if err := json.NewDecoder(resp.Body).Decode(&mpResponse); err != nil {
				http.Error(w, "Erro ao ler resposta do MercadoPago", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mpResponse)
			return

		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}
}

func PagamentoSucesso(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")
		creditsStr := r.URL.Query().Get("credits") // Correção do erro de digitação "creedits" para "credits"
		paymentStatus := r.URL.Query().Get("status")

		fmt.Printf("Pagamento Sucesso: User: %s, Créditos: %s, Status: %s\n", userIDStr, creditsStr, paymentStatus)

		if paymentStatus == "approved" {
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				fmt.Println("Erro ao converter UserID:", err)
				http.Error(w, "Erro interno", http.StatusInternalServerError)
				return
			}

			credits, err := strconv.Atoi(creditsStr)
			if err != nil {
				fmt.Println("Erro ao converter Créditos:", err)
				http.Error(w, "Erro interno", http.StatusInternalServerError) // Correção do erro de digitação "intern" para "interno"
				return
			}

			// Atualizar os créditos do usuario no banco de dados
			_, err = db.Exec("UPDATE users SET credits = credits + ? WHERE id = ?", credits, userID)
			if err != nil {
				fmt.Println("Erro ao atualizar créditos do usuário:", err)
				http.Error(w, "Erro interno", http.StatusInternalServerError)
				return
			}

			fmt.Printf("Créditos (%d) adicionados ao usuário ID %d.\n", credits, userID)

			// Exibir uma mensagem de sucesso para o usuário
			tmpl := template.Must(template.ParseFiles("web/templates/payments/pagamento_sucesso.html")) // Correção do caminho do template
			tmpl.Execute(w, map[string]interface{}{
				"Credits": credits,
			})
			return
		} else {
			// Se o pagamento não foi aprovado (ou outro status), exibir uma mensagem de falha
			tmpl := template.Must(template.ParseFiles("web/templates/payments/pagamento_falhou.html"))
			tmpl.Execute(w, nil)
			return
		}
	}
}

func PagamentoFalhou() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Pagamento falhou...")
		tmpl := template.Must(template.ParseFiles("web/templates/payments/pagamento_falhou.html"))
		tmpl.Execute(w, nil)
	}
}