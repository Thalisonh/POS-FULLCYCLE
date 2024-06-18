package main

import (
	"context"
	"net/http"
	"time"
)

func server() {
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

func Save(ctx context.Context) {

}
