package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const codeLength = 15

// GiftCode representa a estrutura do código de crédito no banco de dados
type GiftCode struct {
	ID             int
	UserID         int
	Code           string
	Credits        int
	PaymentID      string
	PaymentStatus  string
	Redeemed       bool
	RedeemedByUserID sql.NullInt64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// generateUniqueCode gera um código alfanumérico único
func generateUniqueCode() string {
	return uuid.New().String()[:codeLength]
}

// creditUserWithCode gera um código de crédito e o associa ao usuário no banco de dados
func creditUserWithCode(db *sql.DB, userID int, credits int) (string, error) {
	code := generateUniqueCode()
	_, err := db.Exec("INSERT INTO gift_codes (user_id, code, credits, payment_status) VALUES (?, ?, ?, 'pending')", userID, code, credits)
	if err != nil {
		log.Printf("Erro ao inserir código de crédito no banco de dados: %v\n", err)
		return "", fmt.Errorf("erro ao gerar código de crédito")
	}
	log.Printf("Código de crédito gerado para o usuário %d: %s (créditos: %d)\n", userID, code, credits)
	return code, nil
}

// redeemCode verifica e resgata um código de crédito, adicionando os créditos ao usuário
func redeemCode(db *sql.DB, userID int, code string) (int, error) {
	var credits int
	var redeemed bool
	var redeemedByUserIDNullable sql.NullInt64

	row := db.QueryRow("SELECT credits, redeemed, redeemed_by_user_id FROM gift_codes WHERE code = ?", code)
	err := row.Scan(&credits, &redeemed, &redeemedByUserIDNullable)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("código inválido")
	} else if err != nil {
		return 0, fmt.Errorf("erro ao verificar código: %v", err)
	}

	if redeemed {
		if redeemedByUserIDNullable.Valid && int(redeemedByUserIDNullable.Int64) == userID {
			return 0, fmt.Errorf("código já resgatado por você")
		}
		return 0, fmt.Errorf("código já resgatado por outro usuário")
	}

	_, err = db.Exec("UPDATE gift_codes SET redeemed = TRUE, redeemed_by_user_id = ? WHERE code = ?", userID, code)
	if err != nil {
		return 0, fmt.Errorf("erro ao marcar código como resgatado: %v", err)
	}

	_, err = db.Exec("UPDATE users SET credits = credits + ? WHERE id = ?", credits, userID)
	if err != nil {
		return 0, fmt.Errorf("erro ao adicionar créditos ao usuário: %v", err)
	}

	log.Printf("Usuário %d resgatou o código %s, adicionando %d créditos.\n", userID, code, credits)
	return credits, nil
}

// consultarDetalhesPagamentoMP simula a consulta de detalhes do pagamento no Mercado Pago
func consultarDetalhesPagamentoMP(paymentID string) (map[string]interface{}, error) {
	// *** Substitua isso pela sua integração real com a API do Mercado Pago ***
	// Este é apenas um exemplo para simular diferentes status
	if paymentID == "PAYMENT_ID_SIMULADO_APROVADO" {
		return map[string]interface{}{"status": "approved"}, nil
	} else if paymentID == "PAYMENT_ID_SIMULADO_PENDENTE" {
		return map[string]interface{}{"status": "pending"}, nil
	} else if paymentID == "PAYMENT_ID_SIMULADO_REJEITADO" {
		return map[string]interface{}{"status": "rejected"}, nil
	}
	return nil, fmt.Errorf("ID de pagamento simulado inválido")
}

// processPendingPayments verifica periodicamente o status dos pagamentos pendentes
func processPendingPayments(db *sql.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Verificando pagamentos pendentes...")
		rows, err := db.Query("SELECT gc.code, gc.user_id, gc.credits, gc.payment_id FROM gift_codes gc WHERE gc.payment_status = 'pending'")
		if err != nil {
			log.Printf("Erro ao buscar códigos de pagamento pendentes: %v", err)
			continue
		}
		defer rows.Close()

		for rows.Next() {
			var code string
			var userID int
			var credits int
			var paymentID string
			if err := rows.Scan(&code, &userID, &credits, &paymentID); err != nil {
				log.Printf("Erro ao escanear linha de pagamento pendente: %v", err)
				continue
			}

			if paymentID != "" {
				paymentDetails, err := consultarDetalhesPagamentoMP(paymentID)
				if err != nil {
					log.Printf("Erro ao consultar status do pagamento %s: %v", paymentID, err)
					continue
				}

				status, ok := paymentDetails["status"].(string)
				if !ok {
					log.Printf("Status do pagamento %s inválido: %+v", paymentID, paymentDetails)
					continue
				}

				log.Printf("Status do pagamento %s: %s", paymentID, status)

				_, err = db.Exec("UPDATE gift_codes SET payment_status = ? WHERE payment_id = ?", status, paymentID)
				if err != nil {
					log.Printf("Erro ao atualizar status do pagamento %s no banco de dados: %v", paymentID, err)
					continue
				}

				if status == "approved" {
					_, err = db.Exec("UPDATE users SET credits = credits + ? WHERE id = ?", credits, userID)
					if err != nil {
						log.Printf("Erro ao creditar usuário %d após aprovação do pagamento %s: %v", userID, paymentID, err)
						continue
					}
					_, err = db.Exec("UPDATE gift_codes SET redeemed = TRUE WHERE code = ?", code)
					if err != nil {
						log.Printf("Erro ao marcar código %s como resgatado: %v", code, err)
						continue
					}
					log.Printf("Pagamento %s aprovado. Usuário %d creditado com %d créditos (código: %s).", paymentID, userID, credits, code)
				}
			}
		}

		if err := rows.Err(); err != nil {
			log.Printf("Erro ao iterar sobre códigos de pagamento pendentes: %v", err)
		}
	}
}

// GerarCodigoEPagamentoHandler cria o código, simula a criação do pagamento e salva no banco
func GerarCodigoEPagamentoHandler(db *sql.DB) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("user/gerar_codigo.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl.Execute(w, nil)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
			return
		}

		userIDStr := r.FormValue("user_id")
		creditsStr := r.FormValue("credits")

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "ID do usuário inválido", http.StatusBadRequest)
			return
		}

		credits, err := strconv.Atoi(creditsStr)
		if err != nil {
			http.Error(w, "Valor de créditos inválido", http.StatusBadRequest)
			return
		}

		code, err := creditUserWithCode(db, userID, credits)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao gerar código de crédito: %v", err), http.StatusInternalServerError)
			return
		}

		// Simule a geração de um ID de pagamento
		paymentID := "PAYMENT_ID_SIMULADO_PENDENTE_" + code

		_, err = db.Exec("UPDATE gift_codes SET payment_id = ? WHERE code = ?", paymentID, code)
		if err != nil {
			log.Printf("Erro ao atualizar payment_id: %v", err)
		}

		data := map[string]interface{}{
			"Code":      code,
			"PaymentID": paymentID,
		}
		tmpl.Execute(w, data)
	}
}

// ResgatarCreditosHandler é um handler para permitir que o usuário resgate um código
func ResgatarCreditosHandler(db *sql.DB) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/resgatar_codigo.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl.Execute(w, nil)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
			return
		}

		code := r.FormValue("codigo_resgate")
		userIDStr := r.FormValue("user_id") // Assumindo que o ID do usuário está disponível

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "ID do usuário inválido", http.StatusBadRequest)
			return
		}

		credits, err := redeemCode(db, userID, code)
		if err != nil {
			data := map[string]interface{}{
				"Error": err.Error(),
			}
			tmpl.Execute(w, data)
			return
		}

		data := map[string]interface{}{
			"Message": fmt.Sprintf("Código resgatado com sucesso! %d créditos adicionados à sua conta.", credits),
		}
		tmpl.Execute(w, data)
	}
}

// StartPaymentVerification inicia o processo de verificação em background
func StartPaymentVerification(db *sql.DB) {
	go processPendingPayments(db)
}