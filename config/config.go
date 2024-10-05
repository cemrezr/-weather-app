package config

import (
	"github.com/spf13/viper"
	"log"
)

func LoadEnv() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func LoadConfig() *WeatherAppScheme {
	var config WeatherAppScheme

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode configuration into struct: %v", err)
	}
	return &config
}
