package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/rafaelmascaro/Weather-By-CEP-With-Tracing/orchestrator/internal/entity"
	"go.opentelemetry.io/otel/trace"
)

var ErrNotFoundZipcode = errors.New("can not find zipcode")

type LocationClient struct {
	BaseURL      string
	OTELTracer   trace.Tracer
	SpanNameOTEL string
}

type LocationResponse struct {
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}

func NewLocationClient(
	url string,
	tracer trace.Tracer,
	spanName string,
) *LocationClient {
	return &LocationClient{
		BaseURL:      url,
		OTELTracer:   tracer,
		SpanNameOTEL: spanName,
	}
}

func (l *LocationClient) GetLocation(ctx context.Context, cep entity.CEP) (string, error) {
	ctx, span := l.OTELTracer.Start(ctx, l.SpanNameOTEL)
	defer span.End()

	url := strings.ReplaceAll(l.BaseURL, "@CEP", string(cep))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var location LocationResponse
	err = json.Unmarshal(body, &location)
	if err != nil {
		return "", err
	}

	if location.Erro == "true" {
		return "", ErrNotFoundZipcode
	}

	return location.Localidade, nil
}
