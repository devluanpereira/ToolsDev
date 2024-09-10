package models

// ==============================================//
type CnpjData struct {
	Cnpj          string `json:"cnpj"`
	Razao_Social  string `json:"razao_social"`
	Nome_Fantasia string `json:"nome_fantasia"`
	Uf            string `json:"uf"`
	Municipio     string `json:"municipio"`
	Bairro        string `json:"bairro"`
	Logradouro    string `json:"logradouro"`
	Numero        string `json:"numero"`
	Cep           string `json:"cep"`
}

//==============================================//
