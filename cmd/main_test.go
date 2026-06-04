package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AmandaIsrael/cep-weather-api/configs"
	"github.com/stretchr/testify/assert"
)

func newTestConfig() *configs.Config {
	return &configs.Config{
		WeatherAPIURL: "http://api.weatherapi.com/v1/current.json",
		ViaCEPURL:     "http://viacep.com.br/ws/%s/json/",
		Timeout:       5 * time.Second,
		Port:          "8080",
	}
}

func TestSetupServerShouldReturnHTTPHandler(t *testing.T) {
	server := setupServer(newTestConfig())

	assert.NotNil(t, server)
	assert.Implements(t, (*http.Handler)(nil), server)
}

func TestSetupServerShouldRegisterCepRoute(t *testing.T) {
	server := setupServer(newTestConfig())

	req := httptest.NewRequest(http.MethodGet, "/01310100", nil)
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, req)
	assert.NotEqual(t, http.StatusNotFound, recorder.Code)
}
