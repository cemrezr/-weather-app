package config

import (
	"fmt"
	"log"
	"weather-app/internal/repository"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func LoadConfig() *WeatherAppScheme {
	var config WeatherAppScheme

	config.Database.Host = viper.GetString("DB_HOST")
	config.Database.Port = viper.GetInt("DB_PORT")
	config.Database.User = viper.GetString("DB_USER")
	config.Database.Password = viper.GetString("DB_PASSWORD")
	config.Database.DBName = viper.GetString("DB_NAME")
	config.Database.SSLMode = viper.GetString("DB_SSLMODE")
	config.Database.TimeZone = viper.GetString("DB_TIMEZONE")

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode configuration into struct: %v", err)
	}
	return &config
}

func ConnectDatabase(cfg DatabaseConfig) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode, cfg.TimeZone)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("🚀 Connected successfully to the database!")
}

func MigrateDatabase() {
	if err := DB.AutoMigrate(&repository.WeatherQuery{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("🚀 Database migration completed successfully!")
}
