package main

import (
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
