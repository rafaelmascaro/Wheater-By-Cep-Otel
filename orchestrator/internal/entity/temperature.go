package entity

type Temperature struct {
	TempC float64
	TempF float64
	TempK float64
}

func NewTemperature(tempC float64) *Temperature {
	tempF := tempC*1.8 + 32
	tempK := tempC + 273
	return &Temperature{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}
}
