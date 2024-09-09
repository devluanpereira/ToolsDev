package services

import (
	"consulta-cep/internal/models"
	"encoding/json"
	"errors"
	"net/http"
)

// FetchBankData consulta a BrasilAPI para buscar os dados de um ISPB.
func FetchBankData(code string) (*models.BankData, error) {
	resp, err := http.Get("https://brasilapi.com.br/api/banks/v1/" + code)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("ISPB não encontrado")
	}

	var data models.BankData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
