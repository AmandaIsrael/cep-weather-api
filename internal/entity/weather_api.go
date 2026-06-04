package entity

type WeatherAPI struct {
	Current Current `json:"current"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
}
