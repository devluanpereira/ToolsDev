package services

import (
	"consulta-cep/internal/models"
	"encoding/json"
	"errors"
	"net/http"
)

func FetchCnpjData(cnpj string) (*models.CnpjData, error) {
	resp, err := http.Get("https://brasilapi.com.br/api/cnpj/v1/" + cnpj)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("CNPJ não encontrado")
	}

	var data models.CnpjData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
