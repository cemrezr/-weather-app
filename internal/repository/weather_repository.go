// repository/weather_repository.go
package repository

import (
	"gorm.io/gorm"
	"log"
)

type WeatherRepository struct {
	db *gorm.DB
}

func NewWeatherRepository(db *gorm.DB) *WeatherRepository {
	return &WeatherRepository{db: db}
}

func (r *WeatherRepository) CreateWeatherQuery(query *WeatherQuery) error {
	if err := r.db.Create(query).Error; err != nil {
		log.Printf("Failed to insert record into weather_queries for location %s: %v", query.Location, err)
		return err
	}
	log.Printf("Successfully inserted record into weather_queries for location %s with Service1Temp: %.2f, Service2Temp: %.2f", query.Location, query.Service1Temp, query.Service2Temp)
	return nil
}

func (r *WeatherRepository) GetWeatherQueriesByLocation(location string) ([]WeatherQuery, error) {
	var queries []WeatherQuery
	result := r.db.Where("location = ?", location).Find(&queries)
	return queries, result.Error
}

func (r *WeatherRepository) GetAllWeatherQueries() ([]WeatherQuery, error) {
	var queries []WeatherQuery
	result := r.db.Find(&queries)
	return queries, result.Error
}

func (r *WeatherRepository) DeleteWeatherQuery(id uint) error {
	result := r.db.Delete(&WeatherQuery{}, id)
	return result.Error
}
