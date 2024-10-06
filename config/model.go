package config

type WeatherConfig struct {
	ClientSecret string `mapstructure:"WEATHERSERVICE_CLIENT_SECRET"`
	URL          string `mapstructure:"WEATHERSERVICE_URL"`
	Timeout      int    `mapstructure:"WEATHERSERVICE_URL_TIMEOUT"`
}

type WeatherStackConfig struct {
	ClientSecret string `mapstructure:"WEATHERSTACKSERVICE_CLIENT_SECRET"`
	URL          string `mapstructure:"WEATHERSTACKSERVICE_URL"`
	Timeout      int    `mapstructure:"WEATHERSTACKSERVICE_TIMEOUT"`
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}
