package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	fmt.Printf("Start server")
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", GetDolar)
	if err := http.ListenAndServe(":8080", recoverMiddleware(mux)); err != nil {
		fmt.Errorf(err.Error())
	}
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic: %s", r)
				debug.PrintStack()
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
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

	err = Save(ctx, Dolar{response.USDBR.Bid})
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

func connectSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Dolar{})

	return db, nil
}

func Save(ctx context.Context, dolar Dolar) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	db, err := connectSQLite()
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).Create(&dolar).Error
	if err != nil {
		return err
	}

	return nil
}
