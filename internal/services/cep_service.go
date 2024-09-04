package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"consulta-cep/internal/models"
)

// FetchCepData consulta a BrasilAPI para buscar os dados de um CEP.
func FetchCepData(cep string) (*models.CepData, error) {
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v2/" + cep)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("CEP não encontrado")
	}

	var data models.CepData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { // Corrigido o uso do operador ':='
		return nil, err
	}

	return &data, nil
}
