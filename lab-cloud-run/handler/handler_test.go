package handler_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Thalisonh/POS-FULLCYCLE/lab-cloud-run/entities"
	"github.com/Thalisonh/POS-FULLCYCLE/lab-cloud-run/handler"
	"github.com/Thalisonh/POS-FULLCYCLE/lab-cloud-run/services"
	"github.com/stretchr/testify/assert"
)

func TestHandlerService_HandleRequest(t *testing.T) {

	t.Run("Should return error when cep cize is not 8", func(t *testing.T) {
		addressService := &services.IAddressServiceMock{}
		weatherService := &services.IWeatherServiceMock{}
		handlerService := handler.NewHandlerService(addressService, weatherService)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/1234567", nil)

		handlerService.HandleRequest(w, r)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 422, res.StatusCode)
	})

	t.Run("Should return error when fail to get address", func(t *testing.T) {
		addressService := &services.IAddressServiceMock{
			GetAddressError: errors.New("error getting address"),
		}
		weatherService := &services.IWeatherServiceMock{}
		handlerService := handler.NewHandlerService(addressService, weatherService)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/12345678", nil)

		handlerService.HandleRequest(w, r)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 404, res.StatusCode)
	})

	t.Run("Should return error when fail to get weather", func(t *testing.T) {
		addressService := &services.IAddressServiceMock{
			GetAddressResult: entities.Address{
				Zipcode: "12345678",
				City:    "city",
			},
		}
		weatherService := &services.IWeatherServiceMock{
			GetWeatherError: errors.New("error getting weather"),
		}
		handlerService := handler.NewHandlerService(addressService, weatherService)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/12345678", nil)

		handlerService.HandleRequest(w, r)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 404, res.StatusCode)
	})

	t.Run("Should return error when fail to get weather", func(t *testing.T) {
		addressService := &services.IAddressServiceMock{
			GetAddressResult: entities.Address{
				Zipcode: "12345678",
				City:    "city",
			},
		}
		weatherService := &services.IWeatherServiceMock{
			GetWeatherResult: entities.Weather{
				Current: entities.Current{
					Celsius:    10,
					Fahrenheit: 50,
				},
			},
		}
		handlerService := handler.NewHandlerService(addressService, weatherService)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/12345678", nil)

		handlerService.HandleRequest(w, r)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, 200, res.StatusCode)
	})
}
