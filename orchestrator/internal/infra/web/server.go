package web

import (
	"encoding/json"
	"net/http"

	"github.com/rafaelmascaro/weather-api-otel/orchestrator/internal/adapters/api"
	"github.com/rafaelmascaro/weather-api-otel/orchestrator/internal/entity"
	"github.com/rafaelmascaro/weather-api-otel/orchestrator/internal/usecase"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Webserver struct {
	LocationClient  entity.LocationClientInterface
	WeatherClient   entity.WeatherClientInterface
	OTELTracer      trace.Tracer
	RequestNameOTEL string
}

// NewServer creates a new server instance
func NewServer(
	locationClient entity.LocationClientInterface,
	weatherClient entity.WeatherClientInterface,
	otelTracer trace.Tracer,
	requestNameOTEL string,
) *Webserver {
	return &Webserver{
		LocationClient:  locationClient,
		WeatherClient:   weatherClient,
		OTELTracer:      otelTracer,
		RequestNameOTEL: requestNameOTEL,
	}
}

// createServer creates a new server instance with go chi router
func (we *Webserver) CreateServer() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/temp", we.HandleRequest)
	return router
}

func (h *Webserver) HandleRequest(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := h.OTELTracer.Start(ctx, h.RequestNameOTEL)
	defer span.End()

	queryParams := r.URL.Query()
	input := queryParams.Get("CEP")
	getTemp := usecase.NewGetTempUseCase(h.LocationClient, h.WeatherClient)
	output, err := getTemp.Execute(ctx, input)
	if err != nil {
		if err == entity.ErrInvalidZipcode {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		} else if err == api.ErrNotFoundZipcode {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
