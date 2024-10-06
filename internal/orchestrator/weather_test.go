package orchestrator

import (
	"context"
	"fmt"
	"testing"
	"weather-app/internal/repository"
	"weather-app/pkg/weatherclient"
	"weather-app/pkg/weatherstackclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type WeatherClientMock struct {
	mock.Mock
}

func (w *WeatherClientMock) GetWeatherData(ctx context.Context, location string) (*weatherclient.CurrentWeather, error) {
	args := w.Called(ctx, location)
	return args.Get(0).(*weatherclient.CurrentWeather), args.Error(1)
}

type WeatherStackClientMock struct {
	mock.Mock
}

func (w *WeatherStackClientMock) GetWeatherData(ctx context.Context, location string) (*weatherstackclient.CurrentWeather, error) {
	args := w.Called(ctx, location)
	return args.Get(0).(*weatherstackclient.CurrentWeather), args.Error(1)
}

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) CreateWeatherQuery(query *repository.WeatherQuery) error {
	args := r.Called(query)
	return args.Error(0)
}

func (r *RepositoryMock) GetWeatherQueriesByLocation(location string) ([]repository.WeatherQuery, error) {
	args := r.Called(location)
	return args.Get(0).([]repository.WeatherQuery), args.Error(1)
}

func (r *RepositoryMock) GetAllWeatherQueries() ([]repository.WeatherQuery, error) {
	args := r.Called()
	return args.Get(0).([]repository.WeatherQuery), args.Error(1)
}

func (r *RepositoryMock) DeleteWeatherQuery(id uint) error {
	args := r.Called(id)
	return args.Error(0)
}

func TestWeatherOrchestrator_GetAverageTemperaturesBatch(t *testing.T) {
	weatherClient := new(WeatherClientMock)
	weatherStackClient := new(WeatherStackClientMock)
	repo := new(RepositoryMock)

	weatherData := &weatherclient.CurrentWeather{Temperature: 25.0}
	weatherStackData := &weatherstackclient.CurrentWeather{Temperature: 26.0}

	weatherClient.On("GetWeatherData", mock.Anything, "Istanbul").Return(weatherData, nil)
	weatherStackClient.On("GetWeatherData", mock.Anything, "Istanbul").Return(weatherStackData, nil)

	repo.On("CreateWeatherQuery", mock.Anything).Return(nil)

	orchestrator := NewWeatherOrchestrator(weatherClient, weatherStackClient, repo)

	ctx := context.Background()
	locations := []string{"Istanbul"}
	results, err := orchestrator.GetAverageTemperaturesBatch(ctx, locations, 1)

	expectedResult := map[string]float64{"Istanbul": 25.5}

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, results)

	weatherClient.AssertExpectations(t)
	weatherStackClient.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestWeatherOrchestrator_GetAverageTemperaturesBatch_ErrorFromWeatherClient(t *testing.T) {
	weatherClient := new(WeatherClientMock)
	weatherStackClient := new(WeatherStackClientMock)
	repo := new(RepositoryMock)

	weatherClient.On("GetWeatherData", mock.Anything, "Istanbul").Return((*weatherclient.CurrentWeather)(nil), fmt.Errorf("WeatherClient error"))
	weatherStackClient.On("GetWeatherData", mock.Anything, "Istanbul").Return(&weatherstackclient.CurrentWeather{Temperature: 26.0}, nil)

	orchestrator := NewWeatherOrchestrator(weatherClient, weatherStackClient, repo)

	ctx := context.Background()
	locations := []string{"Istanbul"}
	_, err := orchestrator.GetAverageTemperaturesBatch(ctx, locations, 1)

	assert.Error(t, err)

	weatherClient.AssertExpectations(t)
	weatherStackClient.AssertExpectations(t)
}

func TestWeatherOrchestrator_GetAverageTemperaturesBatch_ErrorFromWeatherStackClient(t *testing.T) {
	weatherClient := new(WeatherClientMock)
	weatherStackClient := new(WeatherStackClientMock)
	repo := new(RepositoryMock)

	weatherClient.On("GetWeatherData", mock.Anything, "Istanbul").Return(&weatherclient.CurrentWeather{Temperature: 25.0}, nil)
	weatherStackClient.On("GetWeatherData", mock.Anything, "Istanbul").Return((*weatherstackclient.CurrentWeather)(nil), fmt.Errorf("WeatherStackClient error"))

	orchestrator := NewWeatherOrchestrator(weatherClient, weatherStackClient, repo)

	ctx := context.Background()
	locations := []string{"Istanbul"}
	_, err := orchestrator.GetAverageTemperaturesBatch(ctx, locations, 1)

	assert.Error(t, err)

	weatherClient.AssertExpectations(t)
	weatherStackClient.AssertExpectations(t)
}

func TestWeatherOrchestrator_GetAverageTemperaturesBatch_EmptyLocations(t *testing.T) {
	weatherClient := new(WeatherClientMock)
	weatherStackClient := new(WeatherStackClientMock)
	repo := new(RepositoryMock)

	orchestrator := NewWeatherOrchestrator(weatherClient, weatherStackClient, repo)

	ctx := context.Background()
	var locations []string
	results, err := orchestrator.GetAverageTemperaturesBatch(ctx, locations, 1)

	expectedResult := map[string]float64{}

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, results)

	weatherClient.AssertExpectations(t)
	weatherStackClient.AssertExpectations(t)
	repo.AssertExpectations(t)
}
