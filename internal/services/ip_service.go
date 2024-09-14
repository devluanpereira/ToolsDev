package services

import (
	"consulta-cep/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

const ipApiURL = "http://ipinfo.io/"

func GetIPInfo(ip string) (*models.IPInfo, error) {
	url := fmt.Sprintf("%s%s/json", ipApiURL, ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao realizar a requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou status: %s", resp.Status)
	}

	var ipInfo models.IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&ipInfo); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &ipInfo, nil
}
