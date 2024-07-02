package entities

type Weather struct {
	Day       string
	Condition string // TODO: criar um enum para exibir a descrição refernte a sigla
	Max       int
	Min       int
	IUV       float32
}

type WeatherCoast struct {
	Morning   WeatherCoastData
	Afternoon WeatherCoastData
	Evening   WeatherCoastData
}

type WeatherCoastData struct {
	Day           string
	SeaAgiation   string
	WaveHeight    string
	Direction     string
	WindSpeed     string
	WindDirection string
}
