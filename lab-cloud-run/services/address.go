package services

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Thalisonh/POS-FULLCYCLE/lab-cloud-run/entities"
)

type IAddressService interface {
	GetAddress(cep string) (entities.Address, error)
}

type AddressService struct {
}

func NewAddressService() IAddressService {
	return &AddressService{}
}

func (service *AddressService) GetAddress(cep string) (entities.Address, error) {
	response, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return entities.Address{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return entities.Address{}, err
	}

	var address entities.Address
	json.Unmarshal(body, &address)

	return address, nil
}
