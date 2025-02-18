package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/thalison/POS/lab-observability/configs"
	"github.com/thalison/POS/lab-observability/pkg"
	"github.com/thalison/POS/lab-observability/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
)

var tracer = otel.Tracer("service_a")

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	exporter, err := zipkin.New(config.ZipkinUrl)
	if err != nil {
		fmt.Println(fmt.Sprintf("zipkin exporter: %v", err))
	}
	defer exporter.Shutdown(ctx)

	shutdown, err := pkg.InitProvider("service_a", config.OtelExporterOtlpEndpoint, exporter)
	if err != nil {
		fmt.Errorf("errro init provider")
		panic(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			fmt.Errorf("shutdown")
			panic(err)
		}
	}()

	svc := server.NewServer(&configs.Config{
		OTELTracer: tracer,
	})
	route := svc.CreateServer(config.ServiceNameA)
	handler := NewHandler(&configs.Config{
		ServiceBUrl: config.ServiceBUrl,
	})

	route.Post("/", handler.HandleRequest)

	fmt.Println(fmt.Sprintf("server name %s", config.ServiceNameA))
	fmt.Println(fmt.Sprintf("running on port %s", config.PortServiceA))
	http.ListenAndServe(config.PortServiceA, route)
}

type Handler struct {
	Config *configs.Config
}

func NewHandler(config *configs.Config) *Handler {
	return &Handler{Config: config}
}

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := tracer.Start(ctx, "service_a")
	defer span.End()

	cep := Cep{}
	if err := json.NewDecoder(r.Body).Decode(&cep); err != nil {
		fmt.Errorf(err.Error())

		w.WriteHeader(http.StatusUnprocessableEntity)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid zipcode",
		})

		return
	}

	if len(cep.Cep) != 8 {
		fmt.Errorf("invalid cep")

		w.WriteHeader(http.StatusUnprocessableEntity)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid zipcode",
		})

		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", h.Config.ServiceBUrl, cep.Cep), nil)
	if err != nil {
		fmt.Errorf(err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf(err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Errorf(err.Error())

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
