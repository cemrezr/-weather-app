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
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DBName   string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
	TimeZone string `mapstructure:"DB_TIMEZONE"`
}
