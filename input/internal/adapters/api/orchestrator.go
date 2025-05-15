package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rafaelmascaro/weather-api-otel/input/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Orchestrator struct {
	BaseURL      string
	OTELTracer   trace.Tracer
	SpanNameOTEL string
}

func NewOrchestrator(
	url string,
	tracer trace.Tracer,
	spanName string,
) *Orchestrator {
	return &Orchestrator{
		BaseURL:      url,
		OTELTracer:   tracer,
		SpanNameOTEL: spanName,
	}
}

func (o *Orchestrator) GetTemp(ctx context.Context, cep entity.CEP) (*entity.Temp, error) {
	ctx, span := o.OTELTracer.Start(ctx, o.SpanNameOTEL)
	defer span.End()

	url := strings.ReplaceAll(o.BaseURL, "@CEP", string(cep))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		return nil, entity.ErrInvalidZipcode
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, entity.ErrNotFoundZipcode
	}

	var temp entity.Temp
	err = json.Unmarshal(body, &temp)
	if err != nil {
		return nil, err
	}

	return &temp, nil
}
