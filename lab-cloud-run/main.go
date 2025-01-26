package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/Thalisonh/POS-FULLCYCLE/lab-cloud-run/handler"
	"github.com/Thalisonh/POS-FULLCYCLE/lab-cloud-run/services"
)

func main() {
	// Desabilitar a verificação do certificado SSL
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	addressSvc := services.NewAddressService()
	weatherSvc := services.NewWeatherService()
	handler := handler.NewHandlerService(addressSvc, weatherSvc)

	http.HandleFunc("/{cep}", handler.HandleRequest)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
