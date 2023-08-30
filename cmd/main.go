package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	urlApiCep = "https://cdn.apicep.com/file/apicep/%s.json"
	urlViaCep = "http://viacep.com.br/ws/%s/json/"
)

type Message struct {
	Url     string `json:"url"`
	Error   bool   `json:"error"`
	Content string `json:"content"`
}

func (m *Message) ToString() string {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Informe o CEP!")
		return
	}

	if len(os.Args[1]) != 8 {
		fmt.Println("CEP invalido!")
		return
	}

	channelApiCep := make(chan Message)
	channelViaCep := make(chan Message)

	go doRequest(urlApiCep, os.Args[1], channelApiCep)
	go doRequest(urlViaCep, os.Args[1], channelViaCep)

	select {
	case response := <-channelApiCep:
		fmt.Println(response.ToString())
	case response := <-channelViaCep:
		fmt.Println(response.ToString())
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}

}

func doRequest(urlTemplate, cep string, channel chan Message) {
	url := fmt.Sprintf(urlTemplate, cep)

	resp, err := http.Get(url)
	if err != nil {
		msg := createMessage(url, true, "Erro ao criar a requisicao")
		channel <- msg
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := createMessage(url, true, "Erro ao ler a resposta da requisicao")
		channel <- msg
	}

	channel <- createMessage(url, false, string(body))
}

func createMessage(url string, err bool, content string) Message {
	msg := Message{
		Url:     url,
		Error:   err,
		Content: content,
	}
	return msg
}
