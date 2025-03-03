package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	exporter, err := zipkin.New(config.ZipkinUrl)
	if err != nil {
		fmt.Println(fmt.Sprintf("error initializing zipkin exporter: %v", err))
	}
	defer exporter.Shutdown(ctx)

	shutdown, err := pkg.InitProvider("serviceB", config.OtelExporterOtlpEndpoint, exporter)
	if err != nil {
		fmt.Println(fmt.Sprintf("error initializing provider init: %v", err))
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			fmt.Println(fmt.Sprintf("error initializing provider: %v", err))
		}
	}()

	svc := server.NewServer(&configs.Config{
		OTELTracer: tracer,
	})

	route := svc.CreateServer(config.ServiceNameB)

	route.Post("/{cep}", HandleRequest)

	fmt.Println(fmt.Sprintf("server name %s", config.ServiceNameB))
	fmt.Println(fmt.Sprintf("running on port %s", config.PortB))
	http.ListenAndServe(config.PortB, route)
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
	err = json.Unmarshal(body, &address)
	if err != nil {
		return Address{}, err
	}

	if address.City == "" {
		return Address{}, errors.New("record not found")
	}

	return address, nil
}

func GetWeather(city string) (Weather, error) {
	apiKey := "356249fd69394d598dc213126241511"

	encodedCity := url.QueryEscape(city)

	response, err := http.Get(fmt.Sprintf("https://api.weatherapi.com/v1/current.json?q=%s&key=%s", encodedCity, apiKey))
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
