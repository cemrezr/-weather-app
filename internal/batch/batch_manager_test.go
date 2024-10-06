package batch

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
	"time"
)

type WeatherOrchestratorMock struct {
	mock.Mock
}

func (m *WeatherOrchestratorMock) GetAverageTemperaturesBatch(ctx context.Context, locations []string, requestsCount int) (map[string]float64, error) {
	args := m.Called(ctx, locations, requestsCount)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func TestBatchRequestManager_AddRequest_Success(t *testing.T) {
	orchestratorMock := new(WeatherOrchestratorMock)

	avgTempMap := map[string]float64{"Istanbul": 25.5}
	orchestratorMock.On("GetAverageTemperaturesBatch", mock.Anything, []string{"Istanbul"}, 1).Return(avgTempMap, nil)

	batchManager := NewBatchRequestManager(orchestratorMock)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Timeout süresini uzat
	defer cancel()

	fmt.Println("Starting AddRequest test for Istanbul")

	temp, err := batchManager.AddRequest(ctx, "Istanbul")

	assert.NoError(t, err)
	assert.Equal(t, 25.5, temp)
	orchestratorMock.AssertExpectations(t)
}

func TestBatchRequestManager_AddRequest_ProcessBatchOnTimeout(t *testing.T) {
	orchestratorMock := new(WeatherOrchestratorMock)

	avgTempMap := map[string]float64{"Istanbul": 25.5}
	orchestratorMock.On("GetAverageTemperaturesBatch", mock.Anything, []string{"Istanbul"}, 1).Return(avgTempMap, nil)

	batchManager := NewBatchRequestManager(orchestratorMock)

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	temp, err := batchManager.AddRequest(ctx, "Istanbul")

	assert.NoError(t, err)
	assert.Equal(t, 25.5, temp)
	orchestratorMock.AssertExpectations(t)
}

func TestBatchRequestManager_AddRequest_MaxBatchSize(t *testing.T) {
	orchestratorMock := new(WeatherOrchestratorMock)

	avgTempMap := map[string]float64{"Istanbul": 25.5}
	orchestratorMock.On("GetAverageTemperaturesBatch", mock.Anything, []string{"Istanbul"}, 10).Return(avgTempMap, nil)

	batchManager := NewBatchRequestManager(orchestratorMock)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			temp, err := batchManager.AddRequest(ctx, "Istanbul")
			assert.NoError(t, err)
			assert.Equal(t, 25.5, temp)
		}()
	}
	wg.Wait()
	orchestratorMock.AssertExpectations(t)
}

/*func TestBatchRequestManager_ProcessBatch_Error(t *testing.T) {
  orchestratorMock := new(WeatherOrchestratorMock)

  orchestratorMock.On("GetAverageTemperaturesBatch", mock.Anything, []string{"Istanbul"}, 1).Return(nil, fmt.Errorf("some error"))

  batchManager := NewBatchRequestManager(orchestratorMock)

  ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
  defer cancel()
  _, err := batchManager.AddRequest(ctx, "Istanbul")

  assert.Error(t, err)
  assert.Contains(t, err.Error(), "some error")
  orchestratorMock.AssertExpectations(t)
}*/
