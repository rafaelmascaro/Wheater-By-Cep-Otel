package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

type WeatherClient struct {
	BaseURL      string
	OTELTracer   trace.Tracer
	SpanNameOTEL string
}

type WeatherRequest struct {
	Locations []LocationRequest `json:"locations"`
}

type LocationRequest struct {
	Q string `json:"q"`
}

type WeatherResponse struct {
	Bulk []BulkResponse `json:"bulk"`
}

type BulkResponse struct {
	Query QueryResponse `json:"query"`
}

type QueryResponse struct {
	Current CurrentResponse `json:"current"`
}

type CurrentResponse struct {
	TempC float64 `json:"temp_c"`
}

func NewWeatherClient(
	url string,
	apiKey string,
	tracer trace.Tracer,
	spanName string,
) *WeatherClient {
	baseUrl := strings.ReplaceAll(url, "@APIKEY", apiKey)
	return &WeatherClient{
		BaseURL:      baseUrl,
		OTELTracer:   tracer,
		SpanNameOTEL: spanName,
	}
}

func (w *WeatherClient) GetWeather(ctx context.Context, city string) (float64, error) {
	ctx, span := w.OTELTracer.Start(ctx, w.SpanNameOTEL)
	defer span.End()

	data := WeatherRequest{
		Locations: []LocationRequest{
			{Q: city},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return 0, err
	}

	return weather.Bulk[0].Query.Current.TempC, nil
}
