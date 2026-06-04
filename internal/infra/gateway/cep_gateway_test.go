package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AmandaIsrael/cep-weather-api/configs"
	"github.com/AmandaIsrael/cep-weather-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetViaCEPShouldReturnLocationOnSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"cep":"01310-100","localidade":"São Paulo"}`))
	}))
	defer server.Close()

	gateway := NewCEPGateway(&configs.Config{ViaCEPURL: server.URL + "/%s"})

	result, err := gateway.GetViaCEP(context.Background(), "01310100")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "São Paulo", result.Place)
}

func TestGetViaCEPShouldReturnNotFoundWhenErroIsTrue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"erro": "true"}`))
	}))
	defer server.Close()

	gateway := NewCEPGateway(&configs.Config{ViaCEPURL: server.URL + "/%s"})

	result, err := gateway.GetViaCEP(context.Background(), "99999999")

	assert.ErrorIs(t, err, entity.ErrZipcodeNotFound)
	assert.Nil(t, result)
}

func TestGetViaCEPShouldReturnErrorOnNon200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	gateway := NewCEPGateway(&configs.Config{ViaCEPURL: server.URL + "/%s"})

	result, err := gateway.GetViaCEP(context.Background(), "01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "status 500")
}

func TestGetViaCEPShouldReturnErrorOnInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{ invalid json }"))
	}))
	defer server.Close()

	gateway := NewCEPGateway(&configs.Config{ViaCEPURL: server.URL + "/%s"})

	result, err := gateway.GetViaCEP(context.Background(), "01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetViaCEPShouldReturnErrorOnContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	gateway := NewCEPGateway(&configs.Config{ViaCEPURL: server.URL + "/%s"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	result, err := gateway.GetViaCEP(ctx, "01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetWeatherAPIShouldReturnTemperatureOnSuccess(t *testing.T) {
	var receivedKey, receivedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedKey = r.URL.Query().Get("key")
		receivedQuery = r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer server.Close()

	gateway := NewCEPGateway(&configs.Config{WeatherAPIURL: server.URL, AppKey: "my-key"})

	result, err := gateway.GetWeatherAPI(context.Background(), "São Paulo")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 28.5, result.Current.TempC)
	assert.Equal(t, "my-key", receivedKey)
	assert.Equal(t, "São Paulo", receivedQuery)
}

func TestGetWeatherAPIShouldReturnErrorOnNon200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	gateway := NewCEPGateway(&configs.Config{WeatherAPIURL: server.URL, AppKey: "my-key"})

	result, err := gateway.GetWeatherAPI(context.Background(), "São Paulo")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "status 403")
}

func TestGetWeatherAPIShouldReturnErrorOnInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	gateway := NewCEPGateway(&configs.Config{WeatherAPIURL: server.URL, AppKey: "my-key"})

	result, err := gateway.GetWeatherAPI(context.Background(), "São Paulo")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetWeatherAPIShouldReturnErrorOnContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	gateway := NewCEPGateway(&configs.Config{WeatherAPIURL: server.URL, AppKey: "my-key"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	result, err := gateway.GetWeatherAPI(ctx, "São Paulo")

	assert.Error(t, err)
	assert.Nil(t, result)
}
