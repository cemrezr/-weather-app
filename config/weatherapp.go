package config

var WeatherApp *WeatherAppScheme

type WeatherAppScheme struct {
	Weather      WeatherConfig      `mapstructure:",squash"`
	WeatherStack WeatherStackConfig `mapstructure:",squash"`
	Database     DatabaseConfig     `mapstructure:",squash"`
}
