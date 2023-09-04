package resource

import (
	"fmt"
)

type Printable interface {
	Print()
}

type Message struct {
	ApiName  string    `json:"api_name"`
	Url      string    `json:"url"`
	Response Printable `json:"response"`
}

type ApiCepHttp200Response struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type ApiCepHttp403Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type ApiCepHttp429Response struct {
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	Message    string `json:"message"`
	StatusText string `json:"statusText"`
}

type ViaCepHttp200Response struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        bool   `json:"erro"`
}

type GenericErrorResponse struct {
	Erro    bool   `json:"erro"`
	Message string `json:"message"`
}

func (m Message) Print() {
	fmt.Println("-------------------------------------------------")
	fmt.Println("API: ", m.ApiName)
	fmt.Println("URL: ", m.Url)
	fmt.Println("-------------------------------------------------")
	m.Response.Print()
	fmt.Println("-------------------------------------------------")
}

func (p ApiCepHttp403Response) Print() {
	fmt.Println("A API retornou um erro do tipo HTTP 403")
	fmt.Println("Code: ", p.Code)
	fmt.Println("Message: ", p.Message)
	fmt.Println("Status: ", p.Status)
}

func (p ApiCepHttp429Response) Print() {
	fmt.Println("A API retornou um erro do tipo HTTP 429")
	fmt.Println("OK: ", p.Ok)
	fmt.Println("Message: ", p.Message)
	fmt.Println("Status: ", p.Status)
	fmt.Println("Status Text: ", p.StatusText)
}

func (p ViaCepHttp200Response) Print() {
	if p.Erro {
		fmt.Println("CEP nao localizado")
		return
	}
	fmt.Println("CEP: ", p.Cep)
	fmt.Println("Endereco: ", p.Logradouro)
	fmt.Println("Complemento: ", p.Complemento)
	fmt.Println("Bairro: ", p.Bairro)
	fmt.Println("Cidade: ", p.Localidade)
	fmt.Println("UF: ", p.Uf)
	fmt.Println("DDD: ", p.Ddd)
	fmt.Println("Siafi: ", p.Siafi)
	fmt.Println("IBGE: ", p.Ibge)
	fmt.Println("GIA: ", p.Gia)
}

func (p GenericErrorResponse) Print() {
	fmt.Println("A API retornou um erro generico")
}

func (p ApiCepHttp200Response) Print() {
	fmt.Println("CEP: ", p.Code)
	fmt.Println("Endereco: ", p.Address)
	fmt.Println("Bairro: ", p.District)
	fmt.Println("Cidade: ", p.City)
	fmt.Println("Estado: ", p.State)
}
