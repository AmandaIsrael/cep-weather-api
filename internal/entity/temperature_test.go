package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTemperatureShouldConvertCelsiusToAllScales(t *testing.T) {
	tests := []struct {
		name       string
		celsius    float64
		fahrenheit float64
		kelvin     float64
	}{
		{name: "objective example", celsius: 28.5, fahrenheit: 83.3, kelvin: 301.65},
		{name: "freezing point", celsius: 0, fahrenheit: 32, kelvin: 273.15},
		{name: "boiling point", celsius: 100, fahrenheit: 212, kelvin: 373.15},
		{name: "negative value", celsius: -10, fahrenheit: 14, kelvin: 263.15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temperature := NewTemperature(tt.celsius)

			assert.Equal(t, tt.celsius, temperature.Celsius)
			assert.InDelta(t, tt.fahrenheit, temperature.Fahrenheit, 0.001)
			assert.InDelta(t, tt.kelvin, temperature.Kelvin, 0.001)
		})
	}
}
