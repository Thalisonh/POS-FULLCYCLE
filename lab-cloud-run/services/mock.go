package services

import "github.com/Thalisonh/POS-FULLCYCLE/lab-cloud-run/entities"

type IAddressServiceMock struct {
	GetAddressResult entities.Address
	GetAddressError  error
}

func (mock *IAddressServiceMock) GetAddress(cep string) (entities.Address, error) {
	return mock.GetAddressResult, mock.GetAddressError
}

type IWeatherServiceMock struct {
	GetWeatherResult entities.Weather
	GetWeatherError  error
}

func (mock *IWeatherServiceMock) GetWeather(city string) (entities.Weather, error) {
	return mock.GetWeatherResult, mock.GetWeatherError
}
