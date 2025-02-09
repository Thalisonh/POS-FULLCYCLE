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
	// "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	// "go.opentelemetry.io/otel/propagation"
	// "go.opentelemetry.io/otel/sdk/resource"
	// sdktrace "go.opentelemetry.io/otel/sdk/trace"
	// semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
)

// // todo fazer separado
// func initProvider(serviceName, coolectorURL string) (func(context.Context) error, error) {
// 	ctx := context.Background()

// 	res, err := resource.New(ctx, resource.WithAttributes(
// 		semconv.ServiceName(serviceName),
// 	),
// 	)
// 	if err != nil {
// 		return nil, err // log error
// 	}

// 	ctx, cancel := context.WithTimeout(ctx, time.Second)
// 	defer cancel()
// 	conn, err := grpc.DialContext(ctx, coolectorURL,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithBlock(),
// 	)
// 	if err != nil {
// 		return nil, err // log error
// 	}

// 	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
// 	if err != nil {
// 		return nil, err // log error
// 	}

// 	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
// 	traceProvider := sdktrace.NewTracerProvider(
// 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 		sdktrace.WithResource(res),
// 		sdktrace.WithSpanProcessor(bsp),
// 	)
// 	otel.SetTracerProvider(traceProvider)

// 	otel.SetTextMapPropagator(propagation.TraceContext{})

// 	return traceProvider.Shutdown, nil
// }

func init() {
	viper.AutomaticEnv()
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := pkg.InitProvider("serviceB", viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	tracer := otel.Tracer("serviceB")
	svc := server.NewServer(&server.Config{
		OTELTracer: tracer,
	})

	svc.CreateServer(HandleRequest)

	// http.HandleFunc("/{cep}", HandleRequest)

	// fmt.Println("Server is running on port 8081")
	// http.ListenAndServe(":8081", nil)

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
	cep := r.URL.Path[1:]
	if len(cep) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid zipcode",
		})

		return
	}

	address, err := GetAddress(cep)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "can not find zipcode",
		})

		return
	}

	weather, err := GetWeather(address.City)
	if err != nil {
		fmt.Println(fmt.Sprintf("get weather error: %v", err))

		http.Error(w, "error getting weather", http.StatusNotFound)
		return
	}

	wResponse := &Response{
		City:       address.City,
		Celsius:    weather.Current.Celsius,
		Fahrenheit: weather.Current.Fahrenheit,
		Kelvin:     weather.Current.Celsius + 273,
	}

	json.NewEncoder(w).Encode(wResponse)
}

func GetAddress(cep string) (Address, error) {
	response, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return Address{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Address{}, err
	}

	var address Address
	json.Unmarshal(body, &address)

	return address, nil
}

func GetWeather(city string) (Weather, error) {
	apiKey := "356249fd69394d598dc213126241511"

	response, err := http.Get(fmt.Sprintf("https://api.weatherapi.com/v1/current.json?q=%s&key=%s", city, apiKey))
	if err != nil {
		return Weather{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Weather{}, err
	}

	var weather Weather
	json.Unmarshal(body, &weather)

	return weather, nil
}

type Address struct {
	Zipcode string `json:"cep"`
	City    string `json:"localidade"`
}

type Weather struct {
	Current Current `json:"current"`
}

type Current struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
}

type Response struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}
