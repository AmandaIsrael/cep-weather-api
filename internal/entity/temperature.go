package entity

type Temperature struct {
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
}

func NewTemperature(celsius float64) *Temperature {
	return &Temperature{
		Celsius:    celsius,
		Fahrenheit: celsius*1.8 + 32,
		Kelvin:     celsius + 273.15,
	}
}
