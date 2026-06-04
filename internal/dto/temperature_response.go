package dto

import "github.com/AmandaIsrael/cep-weather-api/internal/entity"

type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func ToTemperatureResponse(t *entity.Temperature) *TemperatureResponse {
	return &TemperatureResponse{
		TempC: t.Celsius,
		TempF: t.Fahrenheit,
		TempK: t.Kelvin,
	}
}
