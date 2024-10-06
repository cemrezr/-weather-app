package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"weather-app/internal/batch"
)

type WeatherHandler struct {
	BatchManager *batch.BatchRequestManager
}

func NewWeatherHandler(batchManager *batch.BatchRequestManager) *WeatherHandler {
	return &WeatherHandler{BatchManager: batchManager}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("q")
	if location == "" {
		http.Error(w, "Location query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	avgTemp, err := h.BatchManager.AddRequest(ctx, location)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch weather data: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"location":    location,
		"temperature": fmt.Sprintf("%.2f", avgTemp),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
