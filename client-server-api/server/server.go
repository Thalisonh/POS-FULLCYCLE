package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	fmt.Printf("Start server")
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", GetDolar)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Errorf(err.Error())
	}
}

func GetDolar(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	response, err := Do(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Println(err)
	}

	resp, err := json.Marshal(Dolar{response.USDBR.Bid})
	if err != nil {
		fmt.Println(err)
	}

	w.Write(resp)
}

func Do(ctx context.Context, method, url string) (*DolarResponse, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	dolar := DolarResponse{}
	err = json.Unmarshal(body, &dolar)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	return &dolar, nil
}

type DolarResponse struct {
	USDBR Response `json:"USDBRL"`
}

type Response struct {
	Bid string `json:"bid"`
}

type Dolar struct {
	Dolar string `json:"dolar"`
}

func Save(ctx context.Context, dolar DolarResponse) error {
	dns := ""

	db, err := gorm.Open(sqlite.Open(dns), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).Save(&dolar.USDBR).Error
	if err != nil {
		return err
	}

	return nil
}
