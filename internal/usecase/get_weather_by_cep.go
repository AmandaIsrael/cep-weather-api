package usecase

import (
	"context"
	"log"

	"github.com/AmandaIsrael/cep-weather-api/internal/dto"
	"github.com/AmandaIsrael/cep-weather-api/internal/entity"
	"github.com/AmandaIsrael/cep-weather-api/internal/infra/gateway"
)

type IGetWeatherByCepUsecase interface {
	Execute(ctx context.Context, cep string) (*dto.TemperatureResponse, error)
}

type GetWeatherByCepUsecase struct {
	gateway gateway.ICEPGateway
}

func NewGetWeatherByCepUsecase(g gateway.ICEPGateway) *GetWeatherByCepUsecase {
	return &GetWeatherByCepUsecase{gateway: g}
}

func (u *GetWeatherByCepUsecase) Execute(ctx context.Context, cep string) (*dto.TemperatureResponse, error) {
	location, err := u.gateway.GetViaCEP(ctx, cep)
	if err != nil {
		if err != entity.ErrZipcodeNotFound {
			log.Printf("[GETWEATHERBYCEP] error fetching location for CEP %s: %v\n", cep, err)
		}
		return nil, err
	}

	weather, err := u.gateway.GetWeatherAPI(ctx, location.Place)
	if err != nil {
		log.Printf("[GETWEATHERBYCEP] error fetching weather for %s: %v\n", location.Place, err)
		return nil, err
	}

	temperature := entity.NewTemperature(weather.Current.TempC)
	return dto.ToTemperatureResponse(temperature), nil
}
