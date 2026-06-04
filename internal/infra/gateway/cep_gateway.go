package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/AmandaIsrael/cep-weather-api/configs"
	"github.com/AmandaIsrael/cep-weather-api/internal/entity"
)

type ICEPGateway interface {
	GetViaCEP(ctx context.Context, cep string) (*entity.ViaCEP, error)
	GetWeatherAPI(ctx context.Context, place string) (*entity.WeatherAPI, error)
}

type CEPGateway struct {
	config *configs.Config
}

func NewCEPGateway(config *configs.Config) *CEPGateway {
	return &CEPGateway{
		config: config,
	}
}

func (c *CEPGateway) GetViaCEP(ctx context.Context, cep string) (*entity.ViaCEP, error) {
	endpoint := fmt.Sprintf(c.config.ViaCEPURL, cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ViaCEP returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var viaCEP entity.ViaCEP
	if err := json.Unmarshal(body, &viaCEP); err != nil {
		return nil, err
	}

	if viaCEP.Erro {
		return nil, entity.ErrZipcodeNotFound
	}

	return &viaCEP, nil
}

func (c *CEPGateway) GetWeatherAPI(ctx context.Context, place string) (*entity.WeatherAPI, error) {
	endpoint, err := url.Parse(c.config.WeatherAPIURL)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("key", c.config.AppKey)
	query.Set("q", place)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WeatherAPI returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather entity.WeatherAPI
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, err
	}

	return &weather, nil
}
