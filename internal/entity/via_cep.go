package entity

import "strings"

type FlexBool bool

func (b *FlexBool) UnmarshalJSON(data []byte) error {
	*b = FlexBool(strings.Trim(string(data), `"`) == "true")
	return nil
}

type ViaCEP struct {
	Cep   string   `json:"cep"`
	Place string   `json:"localidade"`
	Erro  FlexBool `json:"erro"`
}
