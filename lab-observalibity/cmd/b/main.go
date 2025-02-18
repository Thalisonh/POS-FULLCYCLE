package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	// "github.com/openzipkin/zipkin-go"

	"github.com/thalison/POS/lab-observability/configs"
	"github.com/thalison/POS/lab-observability/pkg"
	"github.com/thalison/POS/lab-observability/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
)

var tracer = otel.Tracer("serviceB")

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	exporter, err := zipkin.New("")
	if err != nil {
		fmt.Println(fmt.Sprintf("error initializing zipkin exporter: %v", err))
	}
	defer exporter.Shutdown(ctx)

	shutdown, err := pkg.InitProvider("serviceB", configs.OtelExporterOtlpEndpoint, exporter)
	if err != nil {
		fmt.Println(fmt.Sprintf("error initializing provider init: %v", err))
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			fmt.Println(fmt.Sprintf("error initializing provider: %v", err))
		}
	}()

	svc := server.NewServer(&server.Config{
		OTELTracer: tracer,
	})

	route := svc.CreateServer(configs.ServiceNameB)

	route.Post("/{cep}", HandleRequest)

	fmt.Println(fmt.Sprintf("running on port %s", configs.HostPortServiceB))
	http.ListenAndServe(configs.HostPortServiceB, route)
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	cep := r.URL.Path[1:]
	if len(cep) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid zipcode",
		})

		return
	}

	ctx, getAddressSpan := tracer.Start(ctx, "get address")
	address, err := GetAddress(cep)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "can not find zipcode",
		})

		return
	}
	getAddressSpan.End()

	ctx, getWeatherSpan := tracer.Start(ctx, "get weather")
	weather, err := GetWeather(address.City)
	if err != nil {
		fmt.Println(fmt.Sprintf("get weather error: %v", err))

		http.Error(w, "error getting weather", http.StatusNotFound)
		return
	}
	getWeatherSpan.End()

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
