package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thalison/POS/lab-observability/configs"
)

type WebServer struct {
	Config *configs.Config
}

func NewServer(config *configs.Config) *WebServer {
	return &WebServer{Config: config}
}

func (we *WebServer) CreateServer(serviceName string) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	router.Handle("/metrics", promhttp.Handler())

	return router
}
