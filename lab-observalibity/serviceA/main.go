package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", HandleRequest)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cep := Cep{}
	if err := json.NewDecoder(r.Body).Decode(&cep); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if len(cep.Cep) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid zipcode",
		})

		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:8081/%s", cep.Cep), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

type Cep struct {
	Cep string `json:"cep"`
}

type Response struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}
