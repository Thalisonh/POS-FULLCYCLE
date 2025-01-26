package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Thalisonh/POS-FULLCYCLE/lab-cloud-run/entities"
	"github.com/Thalisonh/POS-FULLCYCLE/lab-cloud-run/services"
)

type IHandler interface {
	HandleRequest(w http.ResponseWriter, r *http.Request)
}

type HandlerService struct {
	addressSvc services.IAddressService
	weatherSvc services.IWeatherService
}

func NewHandlerService(addressSvc services.IAddressService, weatherSvc services.IWeatherService) IHandler {
	return &HandlerService{
		addressSvc: addressSvc,
		weatherSvc: weatherSvc,
	}
}

func (handler *HandlerService) HandleRequest(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Path[1:]
	if len(cep) != 8 {
		fmt.Println(cep)

		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	address, err := handler.addressSvc.GetAddress(cep)
	if err != nil {
		fmt.Println(fmt.Sprintf("get address error: %v", err))

		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	weather, err := handler.weatherSvc.GetWeather(address.City)
	if err != nil {
		fmt.Println(fmt.Sprintf("get weather error: %v", err))

		http.Error(w, "error getting weather", http.StatusNotFound)
		return
	}

	wResponse := &entities.WeatherResponse{
		Celsius:    weather.Current.Celsius,
		Fahrenheit: weather.Current.Fahrenheit,
		Kelvin:     weather.Current.Celsius + 273,
	}

	json.NewEncoder(w).Encode(wResponse)
}
