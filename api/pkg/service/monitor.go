package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tou01/uptime-monitor/pkg/models"
)

// MonitorService handles the monitoring of endpoints
type MonitorService struct {
	endpoints      map[string]*models.Endpoint
	activeCheckers map[string]context.CancelFunc
	client         *http.Client
	mutex          sync.Mutex
	resultChan     chan models.CheckResult
}

// NewMonitorService creates a new monitoring service
func NewMonitorService() *MonitorService {
	return &MonitorService{
		endpoints:      make(map[string]*models.Endpoint),
		activeCheckers: make(map[string]context.CancelFunc),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		resultChan: make(chan models.CheckResult, 100),
	}
}

// AddEndpoint adds a new endpoint to monitor
func (s *MonitorService) AddEndpoint(endpoint models.Endpoint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Generate ID if not provided
	if endpoint.ID == "" {
		endpoint.ID = uuid.New().String()
	}

	// Set created and updated time if not set
	if endpoint.CreatedAt.IsZero() {
		endpoint.CreatedAt = time.Now()
	}
	endpoint.UpdatedAt = time.Now()

	// Store endpoint
	s.endpoints[endpoint.ID] = &endpoint

	// Start monitoring if active
	if endpoint.IsActive {
		s.startMonitoring(endpoint.ID)
	}

	return nil
}

// UpdateEndpoint updates an existing endpoint
func (s *MonitorService) UpdateEndpoint(endpoint models.Endpoint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if endpoint exists
	if _, exists := s.endpoints[endpoint.ID]; !exists {
		return fmt.Errorf("endpoint with ID %s not found", endpoint.ID)
	}

	// Update timestamp
	endpoint.UpdatedAt = time.Now()

	// Stop monitoring if it's active
	if cancel, exists := s.activeCheckers[endpoint.ID]; exists {
		cancel()
		delete(s.activeCheckers, endpoint.ID)
	}

	// Update endpoint
	s.endpoints[endpoint.ID] = &endpoint

	// Restart monitoring if active
	if endpoint.IsActive {
		s.startMonitoring(endpoint.ID)
	}

	return nil
}

// DeleteEndpoint removes an endpoint from monitoring
func (s *MonitorService) DeleteEndpoint(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if endpoint exists
	if _, exists := s.endpoints[id]; !exists {
		return fmt.Errorf("endpoint with ID %s not found", id)
	}

	// Stop monitoring if it's active
	if cancel, exists := s.activeCheckers[id]; exists {
		cancel()
		delete(s.activeCheckers, id)
	}

	// Remove endpoint
	delete(s.endpoints, id)

	return nil
}

// GetEndpoint retrieves an endpoint by ID
func (s *MonitorService) GetEndpoint(id string) (*models.Endpoint, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if endpoint exists
	endpoint, exists := s.endpoints[id]
	if !exists {
		return nil, fmt.Errorf("endpoint with ID %s not found", id)
	}

	return endpoint, nil
}

// GetAllEndpoints retrieves all endpoints for a user
func (s *MonitorService) GetAllEndpoints(userID string) []*models.Endpoint {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var results []*models.Endpoint

	for _, endpoint := range s.endpoints {
		if endpoint.UserID == userID {
			results = append(results, endpoint)
		}
	}

	return results
}

// ListenForResults returns the channel for receiving check results
func (s *MonitorService) ListenForResults() <-chan models.CheckResult {
	return s.resultChan
}

// startMonitoring begins the monitoring process for an endpoint
func (s *MonitorService) startMonitoring(id string) {
	endpoint := s.endpoints[id]

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	s.activeCheckers[id] = cancel

	// Start the monitoring goroutine
	go func() {
		ticker := time.NewTicker(time.Duration(endpoint.Interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				// Monitoring was cancelled
				return
			case <-ticker.C:
				// Time to check the endpoint
				result := s.checkEndpoint(*endpoint)
				s.resultChan <- result
			}
		}
	}()
}

// checkEndpoint performs a single check of an endpoint
func (s *MonitorService) checkEndpoint(endpoint models.Endpoint) models.CheckResult {
	start := time.Now()
	result := models.CheckResult{
		ID:         uuid.New().String(),
		EndpointID: endpoint.ID,
		Timestamp:  time.Now(),
		Success:    false,
	}

	// Create HTTP request
	req, err := http.NewRequest(endpoint.Method, endpoint.URL, bytes.NewBufferString(endpoint.Body))
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// Add headers
	for _, header := range endpoint.Headers {
		req.Header.Add(header.Key, header.Value)
	}

	// Create client with the specified timeout
	client := &http.Client{
		Timeout: time.Duration(endpoint.Timeout) * time.Second,
	}

	// Perform the request
	resp, err := client.Do(req)
	responseTime := time.Since(start).Milliseconds()
	result.ResponseTime = responseTime

	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// Set result fields
	result.Status = resp.StatusCode
	result.StatusText = resp.Status
	result.ResponseBody = string(body)

	// Determine success
	if endpoint.ExpectStatus != 0 {
		result.Success = resp.StatusCode == endpoint.ExpectStatus
	} else {
		// Default success is any 2xx status code
		result.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
	}

	return result
}
