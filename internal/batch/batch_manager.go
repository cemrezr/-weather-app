package batch

import (
	"context"
	"fmt"
	"sync"
	"time"
	"weather-app/internal/orchestrator"
)

type BatchRequest struct {
	Location  string
	Requests  []chan *float64
	Timer     *time.Timer
	Timestamp time.Time
}

type BatchRequestManager struct {
	mu             sync.Mutex
	pendingBatches map[string]*BatchRequest
	Orchestrator   *orchestrator.WeatherOrchestrator
}

func NewBatchRequestManager(orchestrator *orchestrator.WeatherOrchestrator) *BatchRequestManager {
	return &BatchRequestManager{
		pendingBatches: make(map[string]*BatchRequest),
		Orchestrator:   orchestrator,
	}
}

func (m *BatchRequestManager) AddRequest(ctx context.Context, location string) (float64, error) {
	resultCh := make(chan *float64, 1)
	m.mu.Lock()
	batch, exists := m.pendingBatches[location]
	if !exists {
		batch = &BatchRequest{
			Location:  location,
			Requests:  []chan *float64{resultCh},
			Timer:     time.NewTimer(5 * time.Second),
			Timestamp: time.Now(),
		}
		m.pendingBatches[location] = batch
		m.mu.Unlock()

		go func() {
			<-batch.Timer.C
			m.processBatch(ctx, location)
		}()
	} else {
		batch.Requests = append(batch.Requests, resultCh)

		if len(batch.Requests) >= 10 {
			batch.Timer.Stop()
			go m.processBatch(ctx, location)
		}
		m.mu.Unlock()
	}

	select {
	case temp := <-resultCh:
		return *temp, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func (m *BatchRequestManager) processBatch(ctx context.Context, location string) {
	m.mu.Lock()
	batch, exists := m.pendingBatches[location]
	if !exists {
		m.mu.Unlock()
		return
	}
	delete(m.pendingBatches, location)
	m.mu.Unlock()

	fmt.Printf("Processing batch for location: %s with %d requests\n", location, len(batch.Requests))

	avgTempMap, err := m.Orchestrator.GetAverageTemperaturesBatch(ctx, []string{location}, len(batch.Requests))
	if err != nil {
		for _, reqCh := range batch.Requests {
			reqCh <- nil
			close(reqCh)
		}
		fmt.Printf("Failed to get average temperature for location %s: %v\n", location, err)
		return
	}

	avgTemp, ok := avgTempMap[location]
	if !ok {
		fmt.Printf("No temperature data available for location %s\n", location)
		return
	}

	for _, reqCh := range batch.Requests {
		reqCh <- &avgTemp
		close(reqCh)
	}

	fmt.Printf("Successfully processed batch for location: %s with average temperature: %.2f°C\n", location, avgTemp)
}
