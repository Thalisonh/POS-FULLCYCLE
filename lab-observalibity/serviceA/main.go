package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/{cep}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		HandleRequest(w, r)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cep := Cep{}
	// Do something
	if err := json.NewEncoder(w).Encode(&cep); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	// validate size of the cep
	if len(cep.Cep) != 8 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8081", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	// call service b
	w.WriteHeader(http.StatusOK)
}

type Cep struct {
	Cep string `json:"cep"`
}

type Response struct {
}
