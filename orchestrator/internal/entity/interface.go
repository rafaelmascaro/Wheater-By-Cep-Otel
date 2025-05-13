package entity

import "context"

type LocationClientInterface interface {
	GetLocation(context.Context, CEP) (string, error)
}

type WeatherClientInterface interface {
	GetWeather(context.Context, string) (float64, error)
}
