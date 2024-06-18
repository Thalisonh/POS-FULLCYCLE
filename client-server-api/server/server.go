package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Dolar)
	http.ListenAndServe(":8080", mux)
}

func Dolar(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	response := Do(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL")

	w.Write([]byte(response))

}

func Do(ctx context.Context, method, url string) string {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

func Save(ctx context.Context) error {
	dns := ""

	db, err := gorm.Open(sqlite.Open(dns), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).Save(nil).Error
	if err != nil {
		return err
	}

	return nil
}
