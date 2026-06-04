package configs

import (
	"os"
	"time"
)

type Config struct {
	WeatherAPIURL string
	AppKey        string
	ViaCEPURL     string
	Timeout       time.Duration
	Port          string
}

func Load() *Config {
	return &Config{
		WeatherAPIURL: getEnv("WEATHER_API_URL", "http://api.weatherapi.com/v1/current.json"),
		AppKey:        getEnv("APP_KEY", ""),
		ViaCEPURL:     getEnv("VIACEP_URL", "http://viacep.com.br/ws/%s/json/"),
		Timeout:       getDuration("TIMEOUT", 5*time.Second),
		Port:          getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}
