package handler

import (
	"context"
	"fmt"
	"net/http"
	"weather-app/internal/batch"

	"github.com/gorilla/mux"
)

type WeatherHandler struct {
	BatchManager *batch.BatchRequestManager
}

func NewWeatherHandler(batchManager *batch.BatchRequestManager) *WeatherHandler {
	return &WeatherHandler{BatchManager: batchManager}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	if location == "" {
		http.Error(w, "Location is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	avgTemp, err := h.BatchManager.AddRequest(ctx, location)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch weather data: %v", err), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Average Temperature for %s: %.2f°C\n", location, avgTemp)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(response))
	if err != nil {
		fmt.Printf("Failed to write response: %v\n", err)
	}
}
