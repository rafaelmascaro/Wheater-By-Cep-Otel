package usecase

import (
	"context"
	"math"

	"github.com/rafaelmascaro/weather-api-otel/orchestrator/internal/entity"
)

type TempOutput struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type GetTempUseCase struct {
	LocationClient entity.LocationClientInterface
	WeatherClient  entity.WeatherClientInterface
}

func NewGetTempUseCase(
	locationClient entity.LocationClientInterface,
	weatherClient entity.WeatherClientInterface,
) *GetTempUseCase {
	return &GetTempUseCase{
		LocationClient: locationClient,
		WeatherClient:  weatherClient,
	}
}

func (g *GetTempUseCase) Execute(ctx context.Context, input string) (TempOutput, error) {
	cep, err := entity.NewCEP(string(input))
	if err != nil {
		return TempOutput{}, err
	}

	location, err := g.LocationClient.GetLocation(ctx, cep)
	if err != nil {
		return TempOutput{}, err
	}

	tempC, err := g.WeatherClient.GetWeather(ctx, location)
	if err != nil {
		return TempOutput{}, err
	}

	temp := entity.NewTemperature(tempC)

	dto := TempOutput{
		City:  location,
		TempC: math.Round(temp.TempC*10) / 10,
		TempF: math.Round(temp.TempF*10) / 10,
		TempK: math.Round(temp.TempK*10) / 10,
	}

	return dto, nil
}
