package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Thalisonh/POS-FULLCYCLE/lab-cloud-run/entities"
)

type IWeatherService interface {
	GetWeather(city string) (entities.Weather, error)
}

type WeatherService struct {
}

func NewWeatherService() IWeatherService {
	return &WeatherService{}
}

func (service *WeatherService) GetWeather(city string) (entities.Weather, error) {
	apiKey := "356249fd69394d598dc213126241511"

	response, err := http.Get(fmt.Sprintf("https://api.weatherapi.com/v1/current.json?q=%s&key=%s", city, apiKey))
	if err != nil {
		return entities.Weather{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return entities.Weather{}, err
	}

	var weather entities.Weather
	json.Unmarshal(body, &weather)

	return weather, nil
}
