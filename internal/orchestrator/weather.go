package orchestrator

import (
	"context"
	"fmt"
	"log"
	"sync"
	"weather-app/pkg/weatherclient"
	"weather-app/pkg/weatherstackclient"

	"golang.org/x/sync/errgroup"
)

type WeatherOrchestrator struct {
	WeatherClient      *weatherclient.Client
	WeatherStackClient *weatherstackclient.Client
}

func NewWeatherOrchestrator(weatherClient *weatherclient.Client, weatherStackClient *weatherstackclient.Client) WeatherOrchestrator {
	return WeatherOrchestrator{
		WeatherClient:      weatherClient,
		WeatherStackClient: weatherStackClient,
	}
}

func (o *WeatherOrchestrator) GetAverageTemperaturesBatch(ctx context.Context, locations []string) (map[string]float64, error) {
	results := make(map[string]float64)
	mu := &sync.Mutex{}
	g, ctx := errgroup.WithContext(ctx)

	for _, location := range locations {
		loc := location
		g.Go(func() error {
			avgTemp, err := o.getAverageTemperatureForLocation(ctx, loc)
			if err != nil {
				log.Printf("Failed to get average temperature for location %s: %v\n", loc, err)
				return err
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

func (o *WeatherOrchestrator) getAverageTemperatureForLocation(ctx context.Context, location string) (float64, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	temperatures := []float64{}
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		weatherData, err := o.WeatherClient.GetWeatherData(ctx, location)
		if err != nil {
			errChan <- fmt.Errorf("WeatherClient error: %w", err)
			return
		}
		mu.Lock()
		temperatures = append(temperatures, weatherData.Temperature)
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		weatherStackData, err := o.WeatherStackClient.GetWeatherData(ctx, location)
		if err != nil {
			errChan <- fmt.Errorf("WeatherStackClient error: %w", err)
			return
		}
		mu.Lock()
		temperatures = append(temperatures, weatherStackData.Temperature)
		mu.Unlock()
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return 0, err
	}

	var totalTemp float64
	for _, temp := range temperatures {
		totalTemp += temp
	}
	if len(temperatures) == 0 {
		return 0, fmt.Errorf("no temperature data available for location: %s", location)
	}

	avgTemp := totalTemp / float64(len(temperatures))
	return avgTemp, nil
}
