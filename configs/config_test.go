package configs

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadShouldReturnDefaultsWhenEnvIsEmpty(t *testing.T) {
	os.Clearenv()

	config := Load()

	assert.Equal(t, "http://api.weatherapi.com/v1/current.json", config.WeatherAPIURL)
	assert.Equal(t, "", config.AppKey)
	assert.Equal(t, "http://viacep.com.br/ws/%s/json/", config.ViaCEPURL)
	assert.Equal(t, 5*time.Second, config.Timeout)
	assert.Equal(t, "8080", config.Port)
}

func TestLoadShouldOverrideFromEnvVars(t *testing.T) {
	os.Clearenv()
	os.Setenv("WEATHER_API_URL", "https://custom-weather.com/current.json")
	os.Setenv("APP_KEY", "secret-key")
	os.Setenv("VIACEP_URL", "https://custom-viacep.com/%s")
	os.Setenv("TIMEOUT", "10s")
	os.Setenv("PORT", "3000")
	defer os.Clearenv()

	config := Load()

	assert.Equal(t, "https://custom-weather.com/current.json", config.WeatherAPIURL)
	assert.Equal(t, "secret-key", config.AppKey)
	assert.Equal(t, "https://custom-viacep.com/%s", config.ViaCEPURL)
	assert.Equal(t, 10*time.Second, config.Timeout)
	assert.Equal(t, "3000", config.Port)
}

func TestGetEnvShouldReturnDefaultWhenMissing(t *testing.T) {
	os.Clearenv()

	assert.Equal(t, "default", getEnv("NON_EXISTENT_KEY", "default"))
}

func TestGetEnvShouldReturnValueWhenPresent(t *testing.T) {
	os.Setenv("TEST_KEY", "value")
	defer os.Unsetenv("TEST_KEY")

	assert.Equal(t, "value", getEnv("TEST_KEY", "default"))
}

func TestGetDurationShouldReturnDefaultWhenMissing(t *testing.T) {
	os.Clearenv()

	assert.Equal(t, time.Minute, getDuration("NON_EXISTENT_TIMEOUT", time.Minute))
}

func TestGetDurationShouldParseValidValue(t *testing.T) {
	os.Setenv("TEST_TIMEOUT", "10s")
	defer os.Unsetenv("TEST_TIMEOUT")

	assert.Equal(t, 10*time.Second, getDuration("TEST_TIMEOUT", time.Minute))
}

func TestGetDurationShouldFallBackOnInvalidValue(t *testing.T) {
	os.Setenv("TEST_TIMEOUT", "not-a-duration")
	defer os.Unsetenv("TEST_TIMEOUT")

	assert.Equal(t, time.Minute, getDuration("TEST_TIMEOUT", time.Minute))
}
