package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Waelson/internal/resource"
)

const (
	urlApiCep  = "https://cdn.apicep.com/file/apicep/%s.json"
	urlViaCep  = "http://viacep.com.br/ws/%s/json/"
	nameApiCep = "API Cep"
	nameViaCep = "Via Cep"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Informe o CEP!")
		return
	}

	cep := os.Args[1]

	if len(cep) != 8 {
		fmt.Println("CEP invalido!")
		return
	}

	channelApiCep := make(chan resource.Message)
	channelViaCep := make(chan resource.Message)

	go doRequest(urlApiCep, nameApiCep, cep, channelApiCep)
	go doRequest(urlViaCep, nameViaCep, cep, channelViaCep)

	select {
	case response := <-channelApiCep:
		response.Print()
	case response := <-channelViaCep:
		response.Print()
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}

}

func doRequest(urlTemplate, nameApi, cep string, channel chan resource.Message) {
	url := fmt.Sprintf(urlTemplate, cep)

	resp, err := http.Get(url)
	if err != nil {
		genericError := resource.GenericErrorResponse{
			Erro:    true,
			Message: "Ocorreu um erro ao realizar a requisicao",
		}
		msg := createMessage(nameApi, url, genericError)
		channel <- msg
	}

	defer resp.Body.Close()
	result := processResponse(nameApi, resp)
	channel <- createMessage(nameApi, url, result)
}

func processResponse(nameApi string, response *http.Response) resource.Printable {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return resource.GenericErrorResponse{
			Erro:    true,
			Message: "Ocorreu um erro ao processar a resposta",
		}
	}

	if nameApi == nameApiCep {

		if response.StatusCode == 403 {
			var result resource.ApiCepHttp403Response
			err := json.Unmarshal(body, &result)
			if err != nil {
				return resource.GenericErrorResponse{
					Erro:    true,
					Message: "Ocorreu um erro ao desserealizar a resposta",
				}
			}
			return result
		} else if response.StatusCode == 429 {
			var result resource.ApiCepHttp429Response
			err := json.Unmarshal(body, &result)
			if err != nil {
				return resource.GenericErrorResponse{
					Erro:    true,
					Message: "Ocorreu um erro ao desserealizar a resposta",
				}
			}
			return result
		} else {
			var result resource.GenericErrorResponse
			err := json.Unmarshal(body, &result)
			if err != nil {
				return resource.GenericErrorResponse{
					Erro:    true,
					Message: "Ocorreu um erro ao desserializar a resposta",
				}
			}
			return result
		}

	} else if nameApi == nameViaCep {
		if response.StatusCode == 200 {
			var result resource.ViaCepHttp200Response
			err := json.Unmarshal(body, &result)
			if err != nil {
				return resource.GenericErrorResponse{
					Erro:    true,
					Message: "Ocorreu um erro ao desserealizar a resposta",
				}
			}
			return result
		} else {
			var result resource.GenericErrorResponse
			err := json.Unmarshal(body, &result)
			if err != nil {
				return resource.GenericErrorResponse{
					Erro:    true,
					Message: "Ocorreu um erro ao desserializar a resposta",
				}
			}
			return result
		}
	} else {
		panic(fmt.Errorf("nao foi possivel identificar a API de origem da requisicao"))
	}

}

func createMessage(name, url string, response resource.Printable) resource.Message {
	msg := resource.Message{
		ApiName:  name,
		Url:      url,
		Response: response,
	}
	return msg
}
