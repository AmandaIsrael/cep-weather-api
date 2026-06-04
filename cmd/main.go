package main

import (
	"log"
	"net/http"

	"github.com/AmandaIsrael/cep-weather-api/configs"
	"github.com/AmandaIsrael/cep-weather-api/internal/infra/gateway"
	"github.com/AmandaIsrael/cep-weather-api/internal/infra/web/handlers"
	"github.com/AmandaIsrael/cep-weather-api/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config := configs.Load()
	server := setupServer(config)

	log.Printf("server listening on :%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, server); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func setupServer(config *configs.Config) http.Handler {
	cepGateway := gateway.NewCEPGateway(config)
	weatherUsecase := usecase.NewGetWeatherByCepUsecase(cepGateway)
	cepHandler := handlers.NewCepHandler(weatherUsecase, config)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/{cep}", cepHandler.GetCEP)

	return r
}
