package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AmandaIsrael/cep-weather-api/configs"
	"github.com/AmandaIsrael/cep-weather-api/internal/entity"
	"github.com/AmandaIsrael/cep-weather-api/internal/usecase"
	"github.com/AmandaIsrael/cep-weather-api/pkg"
	"github.com/go-chi/chi/v5"
)

type CepHandler struct {
	usecase usecase.IGetWeatherByCepUsecase
	config  *configs.Config
}

func NewCepHandler(uc usecase.IGetWeatherByCepUsecase, config *configs.Config) *CepHandler {
	return &CepHandler{
		usecase: uc,
		config:  config,
	}
}

func (h *CepHandler) GetCEP(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	if !pkg.IsValidCEP(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.config.Timeout)
	defer cancel()

	result, err := h.usecase.Execute(ctx, cep)
	if err != nil {
		if errors.Is(err, entity.ErrZipcodeNotFound) {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
