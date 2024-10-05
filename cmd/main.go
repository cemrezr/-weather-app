package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"weather-app/config"
	"weather-app/internal/batch"
	"weather-app/internal/handler"
	"weather-app/internal/orchestrator"
	"weather-app/pkg/weatherclient"
	"weather-app/pkg/weatherstackclient"

	"github.com/gorilla/mux"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("🚀 Application is starting...")

	config.LoadEnv()
	appConfig := config.LoadConfig()
	log.Printf("Loaded Configuration: %+v\n", appConfig)

	weatherClient := weatherclient.NewClient(weatherclient.Config{
		BaseURL: appConfig.Weather.URL,
		APIKey:  appConfig.Weather.ClientSecret,
		Timeout: time.Duration(appConfig.Weather.Timeout) * time.Second,
	})

	weatherStackClient := weatherstackclient.NewClient(weatherstackclient.Config{
		BaseURL: appConfig.WeatherStack.URL,
		APIKey:  appConfig.WeatherStack.ClientSecret,
		Timeout: time.Duration(appConfig.WeatherStack.Timeout) * time.Second,
	})

	weatherOrchestrator := orchestrator.NewWeatherOrchestrator(weatherClient, weatherStackClient)

	batchManager := batch.NewBatchRequestManager(&weatherOrchestrator)

	weatherHandler := handler.NewWeatherHandler(batchManager)

	r := mux.NewRouter()
	r.HandleFunc("/weather/{location}", weatherHandler.GetWeather).Methods(http.MethodGet)

	log.Println("🌍 Server is starting at port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
