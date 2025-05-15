package web

import (
	"encoding/json"
	"net/http"

	"github.com/rafaelmascaro/weather-api-otel/input/internal/entity"
	"github.com/rafaelmascaro/weather-api-otel/input/internal/usecase"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Webserver struct {
	Orchestrator    entity.OrchestratorInterface
	OTELTracer      trace.Tracer
	RequestNameOTEL string
}

// NewServer creates a new server instance
func NewServer(
	Orchestrator entity.OrchestratorInterface,
	otelTracer trace.Tracer,
	requestNameOTEL string,
) *Webserver {
	return &Webserver{
		Orchestrator:    Orchestrator,
		OTELTracer:      otelTracer,
		RequestNameOTEL: requestNameOTEL,
	}
}

// createServer creates a new server instance with go chi router
func (we *Webserver) CreateServer() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/temp", we.HandleRequest)
	return router
}

func (h *Webserver) HandleRequest(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := h.OTELTracer.Start(ctx, h.RequestNameOTEL)
	defer span.End()

	var input usecase.TempInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getTemp := usecase.NewGetTempUseCase(h.Orchestrator)
	output, err := getTemp.Execute(ctx, input)

	if err != nil {
		if err == entity.ErrInvalidZipcode {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		} else if err == entity.ErrNotFoundZipcode {
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
