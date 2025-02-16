package server

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
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

func (we *WebServer) CreateServer() *chi.Mux {
	router := chi.NewRouter()

	tracer, err := newTracer()
	if err != nil {
		log.Fatal(err)
	}

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60))
	router.Use(zipkinhttp.NewServerMiddleware(
		tracer,
		zipkinhttp.SpanName("request")))

	router.Handle("/metrics", promhttp.Handler())

	return router
}

const endpointURL = "http://localhost:9411/api/v2/spans"

func newTracer() (*zipkin.Tracer, error) {
	// The reporter sends traces to zipkin server
	reporter := reporterhttp.NewReporter(endpointURL)

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: "my_service", Port: 8080}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	t, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, err
	}

	return t, err
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
}
