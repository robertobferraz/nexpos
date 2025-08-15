package dto

type GetCepResponse struct {
	Cep         *string `json:"cep"`
	Logradouro  *string `json:"logradouro"`
	Complemento *string `json:"complemento"`
	Bairro      *string `json:"bairro"`
	Localidade  *string `json:"localidade"`
	Uf          *string `json:"uf"`
	Ibge        *string `json:"ibge"`
}
