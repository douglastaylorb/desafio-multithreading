package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type EnderecoViaCep struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
}

type ResultadoViaCep struct {
	Data EnderecoViaCep `json:"data"`
	Err  error          `json:"error"`
}

type EnderecoBrasilAPI struct {
	CEP          string `json:"cep"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}
type ResultadoBrasilAPI struct {
	Data EnderecoBrasilAPI `json:"data"`
	Err  error             `json:"error"`
}

func main() {
	cep := "01153000"
	ch1 := make(chan ResultadoBrasilAPI)
	ch2 := make(chan ResultadoViaCep)

	go fetchBrasilAPI(cep, ch1)
	go fetchViaCEP(cep, ch2)

	select {
	case res := <-ch1:
		fmt.Printf("BrasilAPI foi mais rápida!\n")
		fmt.Printf("Endereço: %s, %s, %s, %s, %s\n", res.Data.CEP, res.Data.Street, res.Data.Neighborhood, res.Data.City, res.Data.State)
	case res := <-ch2:
		fmt.Printf("ViaCEP foi mais rápida!\n")
		fmt.Printf("Endereço: %s, %s, %s, %s, %s\n", res.Data.CEP, res.Data.Logradouro, res.Data.Bairro, res.Data.Localidade, res.Data.UF)
	case <-time.After(1 * time.Second):
		fmt.Printf("Timeout: Nenhuma resposta recebida dentro do prazo.\n")
	}
}

func fetchBrasilAPI(cep string, ch chan<- ResultadoBrasilAPI) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/" + cep)
	resp, err := http.Get(url)
	if err != nil {
		ch <- ResultadoBrasilAPI{Err: err}
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data EnderecoBrasilAPI
	err = json.Unmarshal(body, &data)

	// time.Sleep(2 * time.Second)
	ch <- ResultadoBrasilAPI{Data: data}
}

func fetchViaCEP(cep string, ch chan<- ResultadoViaCep) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		ch <- ResultadoViaCep{Err: err}
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data EnderecoViaCep
	err = json.Unmarshal(body, &data)

	// time.Sleep(2 * time.Second)
	ch <- ResultadoViaCep{Data: data}
}
