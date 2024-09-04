package models

type CnpjData struct {
	Cnpj         string `json:"cnpj"`
	RazaoSocial  string `json:"razaosocial"`
	NomeFantasia string `json:"nomefantasia"`
	Uf           string `json:"uf"`
	Municipio    string `json:"municipio"`
	Bairro       string `json:"bairro"`
	Logradouro   string `json:"logradouro"`
	Numero       string `json:"numero"`
	Cep          string `json:"cep"`
}
