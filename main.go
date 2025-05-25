package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BrazilCEP struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func getCepBrazilApi1(ch chan<- BrazilCEP, chErr chan<- error) {
	url := fmt.Sprintf(`https://brasilapi.com.br/api/cep/v1/01153000`)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		chErr <- errors.New(fmt.Sprintf("fail to create request: %s", err))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		chErr <- errors.New(fmt.Sprintf("fail to get response: %s", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		chErr <- errors.New(fmt.Sprintf("fail to read body: %s", err))
	}

	var cep BrazilCEP
	err = json.Unmarshal(body, &cep)
	if err != nil {
		chErr <- errors.New(fmt.Sprintf("fail to unmarshal body: %s", err))
	}

	ch <- cep
	chErr <- nil
}

func getCepViaApi2(ch chan<- ViaCEP, chErr chan<- error) {
	url := fmt.Sprintf(`http://viacep.com.br/ws/01153000/json/`)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		chErr <- errors.New(fmt.Sprintf("fail to create request: %s", err))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		chErr <- errors.New(fmt.Sprintf("fail to get response: %s", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		chErr <- errors.New(fmt.Sprintf("fail to read body: %s", err))
	}

	var cep ViaCEP
	err = json.Unmarshal(body, &cep)
	if err != nil {
		chErr <- errors.New(fmt.Sprintf("fail to unmarshal body: %s", err))
	}

	ch <- cep
	chErr <- nil
}

func main() {
	chBrazilCep := make(chan BrazilCEP)
	chViaCep := make(chan ViaCEP)
	chError := make(chan error)

	go getCepBrazilApi1(chBrazilCep, chError)
	go getCepViaApi2(chViaCep, chError)

	select {
	case cep1 := <-chBrazilCep:
		fmt.Println("cep1: ", cep1)
	case cep2 := <-chViaCep:
		fmt.Println("cep2: ", cep2)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout")
	case err := <-chError:
		fmt.Println("error: ", err)
	}
}
