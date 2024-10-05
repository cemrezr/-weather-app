package weatherclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Config struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		baseURL:    cfg.BaseURL,
		apiKey:     cfg.APIKey,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

type CurrentWeather struct {
	Temperature float64 `json:"temp_c"`
}

type WeatherResponse struct {
	Current CurrentWeather `json:"current"`
}

func (c *Client) GetWeatherData(ctx context.Context, location string) (*CurrentWeather, error) {
	url := fmt.Sprintf("%s/forecast.json?key=%s&q=%s&days=1&aqi=no&alerts=no", c.baseURL, c.apiKey, location)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Weather API Response Body: %s\n", string(body))

	var result WeatherResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Failed to decode response body: %v\n", err)
		log.Printf("Body content: %s\n", string(body))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Decoded Weather Response: %+v\n", result)

	return &result.Current, nil
}
