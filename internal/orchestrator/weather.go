package orchestrator

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"weather-app/internal/repository"
	"weather-app/pkg/weatherclient"
	"weather-app/pkg/weatherstackclient"

	"golang.org/x/sync/errgroup"
)

type WeatherOrchestrator struct {
	WeatherClient      *weatherclient.Client
	WeatherStackClient *weatherstackclient.Client
	Repository         *repository.WeatherRepository
}

// NewWeatherOrchestrator, WeatherOrchestrator oluşturur.
func NewWeatherOrchestrator(weatherClient *weatherclient.Client, weatherStackClient *weatherstackclient.Client, repo *repository.WeatherRepository) *WeatherOrchestrator {
	return &WeatherOrchestrator{
		WeatherClient:      weatherClient,
		WeatherStackClient: weatherStackClient,
		Repository:         repo,
	}
}

func (o *WeatherOrchestrator) GetAverageTemperaturesBatch(ctx context.Context, locations []string, requestsCount int) (map[string]float64, error) {
	results := make(map[string]float64)
	mu := &sync.Mutex{}
	g, ctx := errgroup.WithContext(ctx)

	for _, location := range locations {
		loc := location
		g.Go(func() error {
			service1Temp, service2Temp, err := o.getTemperaturesForLocation(ctx, loc)
			if err != nil {
				log.Printf("Failed to get temperatures for location %s: %v\n", loc, err)
				return err
			}

			avgTemp := CalculateAverageTemperature(service1Temp, service2Temp)

			query := &repository.WeatherQuery{
				Location:     loc,
				Service1Temp: service1Temp,
				Service2Temp: service2Temp,
				RequestCount: requestsCount,
				CreatedAt:    time.Now(),
			}

			if err := o.Repository.CreateWeatherQuery(query); err != nil {
				log.Printf("Failed to insert weather query for location %s: %v\n", loc, err)
			}

			mu.Lock()
			results[loc] = avgTemp
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}

func (o *WeatherOrchestrator) getTemperaturesForLocation(ctx context.Context, location string) (float64, float64, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	temperatures := make([]float64, 2) // İlk eleman service1Temp, ikinci eleman service2Temp olacak
	errChan := make(chan error, 2)

	// WeatherClient'ten veri çekme
	wg.Add(1)
	go func() {
		defer wg.Done()
		weatherData, err := o.WeatherClient.GetWeatherData(ctx, location)
		if err != nil {
			errChan <- fmt.Errorf("WeatherClient error: %w", err)
			return
		}
		mu.Lock()
		temperatures[0] = weatherData.Temperature
		mu.Unlock()
	}()

	// WeatherStackClient'ten veri çekme
	wg.Add(1)
	go func() {
		defer wg.Done()
		weatherStackData, err := o.WeatherStackClient.GetWeatherData(ctx, location)
		if err != nil {
			errChan <- fmt.Errorf("WeatherStackClient error: %w", err)
			return
		}
		mu.Lock()
		temperatures[1] = weatherStackData.Temperature
		mu.Unlock()
	}()

	// Tüm goroutine'lerin tamamlanmasını bekle
	wg.Wait()
	close(errChan)

	// Hata varsa döndür
	for err := range errChan {
		return 0, 0, err
	}

	return temperatures[0], temperatures[1], nil
}

func CalculateAverageTemperature(temp1, temp2 float64) float64 {
	return (temp1 + temp2) / 2.0
}
