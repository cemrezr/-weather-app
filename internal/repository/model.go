package repository

import "time"

type WeatherQuery struct {
	ID           uint      `gorm:"primaryKey"`
	Location     string    `gorm:"type:text"`
	Service1Temp float64   `gorm:"type:float"`
	Service2Temp float64   `gorm:"type:float"`
	RequestCount int       `gorm:"type:int"`
	CreatedAt    time.Time `gorm:"type:timestamp"`
}
