package server

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

func (we *WebServer) CreateServer(handlerRequest func(w http.ResponseWriter, r *http.Request)) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60))

	router.Handle("/metrics", promhttp.Handler())
	router.Get("/", handlerRequest)

	return router
}

func (h *WebServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, spanStart := h.Config.OTELTracer.Start(ctx, "Span Start")
	defer spanStart.End()

	ctx, span := h.Config.OTELTracer.Start(ctx, h.Config.RequestNameOTEL)
	defer span.End()

	var req *http.Request
	var err error

	req, err = http.NewRequestWithContext(ctx, h.Config.ExternalCallMethod, h.Config.ExternalCallUrl, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	h.Config.Content = string(bodyBytes) // response

	return
	// }
}
