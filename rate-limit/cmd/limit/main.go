package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/configs"
	"github.com/Thalisonh/POS-FULLCYCLE/rate-limit/internal/server"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	server.StartServer(config)
}

func foo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func middleware(f http.HandlerFunc, config *configs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("API_KEY")
		ip, _ := getIP(r)
		log.Println(apiKey, ip)

		if apiKey == "" {
			http.Error(w, "API Key is required", http.StatusUnauthorized)
			return
		}
		f(w, r)
	}
}

func getIP(r *http.Request) (string, error) {
	// Verifica o cabeçalho X-Forwarded-For
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0]), nil
	}

	// Verifica o cabeçalho X-Real-IP
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP, nil
	}

	// Usa o endereço remoto como último recurso
	ip := strings.Split(r.RemoteAddr, ":")[0]
	return ip, nil
}

func validateJWT() bool {
	return true
}
