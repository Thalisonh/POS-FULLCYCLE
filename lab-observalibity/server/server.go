package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	OTELTracer         trace.Tracer
	RequestNameOTEL    string
	ExternalCallUrl    string
	ExternalCallMethod string
	Content            string
}

type WebServer struct {
	Config *Config
}

func NewServer(config *Config) *WebServer {
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

const endpointURL = "http://localhost:9411/api/v2/spans"
