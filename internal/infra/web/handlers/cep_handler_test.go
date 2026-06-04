package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AmandaIsrael/cep-weather-api/configs"
	"github.com/AmandaIsrael/cep-weather-api/internal/dto"
	"github.com/AmandaIsrael/cep-weather-api/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUsecase struct {
	mock.Mock
}

func (m *mockUsecase) Execute(ctx context.Context, cep string) (*dto.TemperatureResponse, error) {
	args := m.Called(ctx, cep)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TemperatureResponse), args.Error(1)
}

func newRequest(cep string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/"+cep, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cep", cep)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func newHandler(uc *mockUsecase) *CepHandler {
	return NewCepHandler(uc, &configs.Config{Timeout: time.Second})
}

func TestGetCEPShouldReturn200WithTemperaturesOnSuccess(t *testing.T) {
	uc := new(mockUsecase)
	uc.On("Execute", mock.Anything, "01310100").
		Return(&dto.TemperatureResponse{TempC: 28.5, TempF: 83.3, TempK: 301.65}, nil)

	handler := newHandler(uc)
	recorder := httptest.NewRecorder()
	handler.GetCEP(recorder, newRequest("01310100"))

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var body dto.TemperatureResponse
	assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	assert.Equal(t, 28.5, body.TempC)
	assert.Equal(t, 83.3, body.TempF)
	assert.Equal(t, 301.65, body.TempK)
	uc.AssertExpectations(t)
}

func TestGetCEPShouldReturn422OnInvalidZipcode(t *testing.T) {
	uc := new(mockUsecase)
	handler := newHandler(uc)

	invalidCEPs := []string{"123", "12345678a", "123456789", "abcdefgh", "01310-100", ""}
	for _, cep := range invalidCEPs {
		t.Run(cep, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler.GetCEP(recorder, newRequest(cep))

			assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			assert.Contains(t, recorder.Body.String(), "invalid zipcode")
		})
	}

	uc.AssertNotCalled(t, "Execute")
}

func TestGetCEPShouldReturn404WhenZipcodeNotFound(t *testing.T) {
	uc := new(mockUsecase)
	uc.On("Execute", mock.Anything, "99999999").
		Return(nil, entity.ErrZipcodeNotFound)

	handler := newHandler(uc)
	recorder := httptest.NewRecorder()
	handler.GetCEP(recorder, newRequest("99999999"))

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "can not find zipcode")
	uc.AssertExpectations(t)
}

func TestGetCEPShouldReturn500OnInternalError(t *testing.T) {
	uc := new(mockUsecase)
	uc.On("Execute", mock.Anything, "01310100").
		Return(nil, errors.New("weather api unavailable"))

	handler := newHandler(uc)
	recorder := httptest.NewRecorder()
	handler.GetCEP(recorder, newRequest("01310100"))

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	uc.AssertExpectations(t)
}
