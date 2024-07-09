package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCepResponse struct {
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

type BrasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {
	ctx := context.Background()

	cep := "14405415"

	viaCep := make(chan string)
	brasilApi := make(chan string)

	go func() {
		payload := BrasilApiResponse{}
		brasilResp, err := request(ctx, "https://brasilapi.com.br/api/cep/v1/"+cep, "GET", &payload)
		if err != nil {
			fmt.Println(err)
		}
		brasilApi <- brasilResp
	}()

	go func() {
		payload := ViaCepResponse{}
		viaCepResp, err := request(ctx, "http://viacep.com.br/ws/"+cep+"/json/", "GET", &payload)
		if err != nil {
			fmt.Println(err)
		}
		viaCep <- viaCepResp
	}()

	select {
	case b := <-brasilApi:
		fmt.Printf("%v", b)
	case v := <-viaCep:
		fmt.Printf("%v", v)
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}
}

func request(ctx context.Context, url, method string, payload any) (string, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
