package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"
	"github.com/thalison/POS/lab-observability/pkg"
	"github.com/thalison/POS/lab-observability/server"
	"go.opentelemetry.io/otel"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := pkg.InitProvider("serviceA", viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	tracer := otel.Tracer("serviceA")
	svc := server.NewServer(&server.Config{
		OTELTracer: tracer,
	})

	svc.CreateServer(HandleRequest)

	select {
	case <-sigCh:
		log.Println("shutting down")
	case <-ctx.Done():
		log.Println("context done")
	}

	_, shotdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shotdownCancel()
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	// carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	// ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	// ctx, spanStart := h.Config.OTELTracer.Start(ctx, "Span Start")
	// defer spanStart.End()

	// ctx, span := h.Config.OTELTracer.Start(ctx, h.Config.RequestNameOTEL)
	// defer span.End()

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
