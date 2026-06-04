package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/AmandaIsrael/cep-weather-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCEPGateway struct {
	mock.Mock
}

func (m *mockCEPGateway) GetViaCEP(ctx context.Context, cep string) (*entity.ViaCEP, error) {
	args := m.Called(ctx, cep)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ViaCEP), args.Error(1)
}

func (m *mockCEPGateway) GetWeatherAPI(ctx context.Context, place string) (*entity.WeatherAPI, error) {
	args := m.Called(ctx, place)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WeatherAPI), args.Error(1)
}

func TestExecuteShouldReturnConvertedTemperaturesOnSuccess(t *testing.T) {
	gateway := new(mockCEPGateway)
	gateway.On("GetViaCEP", mock.Anything, "01310100").
		Return(&entity.ViaCEP{Place: "São Paulo"}, nil)
	gateway.On("GetWeatherAPI", mock.Anything, "São Paulo").
		Return(&entity.WeatherAPI{Current: entity.Current{TempC: 28.5}}, nil)

	uc := NewGetWeatherByCepUsecase(gateway)
	result, err := uc.Execute(context.Background(), "01310100")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 28.5, result.TempC)
	assert.InDelta(t, 83.3, result.TempF, 0.001)
	assert.InDelta(t, 301.65, result.TempK, 0.001)
	gateway.AssertExpectations(t)
}

func TestExecuteShouldReturnNotFoundWhenViaCEPDoesNotFindZipcode(t *testing.T) {
	gateway := new(mockCEPGateway)
	gateway.On("GetViaCEP", mock.Anything, "99999999").
		Return(nil, entity.ErrZipcodeNotFound)

	uc := NewGetWeatherByCepUsecase(gateway)
	result, err := uc.Execute(context.Background(), "99999999")

	assert.ErrorIs(t, err, entity.ErrZipcodeNotFound)
	assert.Nil(t, result)
	gateway.AssertNotCalled(t, "GetWeatherAPI")
}

func TestExecuteShouldPropagateViaCEPInfraError(t *testing.T) {
	gateway := new(mockCEPGateway)
	gateway.On("GetViaCEP", mock.Anything, "01310100").
		Return(nil, errors.New("connection refused"))

	uc := NewGetWeatherByCepUsecase(gateway)
	result, err := uc.Execute(context.Background(), "01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
	gateway.AssertNotCalled(t, "GetWeatherAPI")
}

func TestExecuteShouldPropagateWeatherAPIError(t *testing.T) {
	gateway := new(mockCEPGateway)
	gateway.On("GetViaCEP", mock.Anything, "01310100").
		Return(&entity.ViaCEP{Place: "São Paulo"}, nil)
	gateway.On("GetWeatherAPI", mock.Anything, "São Paulo").
		Return(nil, errors.New("weather api unavailable"))

	uc := NewGetWeatherByCepUsecase(gateway)
	result, err := uc.Execute(context.Background(), "01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
	gateway.AssertExpectations(t)
}
