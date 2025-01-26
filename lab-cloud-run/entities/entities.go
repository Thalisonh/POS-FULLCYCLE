package entities

type Address struct {
	Zipcode string `json:"cep"`
	City    string `json:"localidade"`
}

type Weather struct {
	Current Current `json:"current"`
}

type Current struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
}

type WeatherResponse struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}
