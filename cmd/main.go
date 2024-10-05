package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"weather-app/config"
	"weather-app/pkg/weatherclient"
	"weather-app/pkg/weatherstackclient"

	"github.com/gorilla/mux"
)

func main() {
	log.SetOutput(os.Stdout)

	log.Println("Application is starting...")

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

	r := mux.NewRouter()
	r.HandleFunc("/weather/{location}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		location := vars["location"]

		log.Printf("Fetching weather data for location: %s", location)

		weatherData, err := weatherClient.GetWeatherData(context.Background(), location)
		if err != nil {
			log.Printf("Failed to fetch weather data from Weather API: %v\n", err)
			http.Error(w, fmt.Sprintf("Failed to fetch weather data from Weather API: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Weather API Temperature for %s: %.2f°C\n", location, weatherData.Temperature)

		weatherStackData, err := weatherStackClient.GetWeatherData(context.Background(), location)
		if err != nil {
			log.Printf("Failed to fetch weather data from WeatherStack API: %v\n", err)
			http.Error(w, fmt.Sprintf("Failed to fetch weather data from WeatherStack API: %v", err), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf(
			"Weather API Temperature: %.2f°C\nWeatherStack Temperature: %.2f°C\n",
			weatherData.Temperature,
			weatherStackData.Temperature,
		)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Printf("Failed to write response: %v", err)
		}

		log.Printf("Response successfully sent for location: %s", location)
	}).Methods(http.MethodGet)

	log.Println("Server is starting at port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
