package repository

import (
	"time"

	"github.com/tou01/uptime-monitor/pkg/models"
)

// Repository defines the interface for data storage
type Repository interface {
	// User repository methods
	CreateUser(user models.User) error
	GetUserByID(id string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	UpdateUser(user models.User) error
	DeleteUser(id string) error

	// Endpoint repository methods
	CreateEndpoint(endpoint models.Endpoint) error
	GetEndpointByID(id string) (models.Endpoint, error)
	GetEndpointsByUserID(userID string) ([]models.Endpoint, error)
	UpdateEndpoint(endpoint models.Endpoint) error
	DeleteEndpoint(id string) error

	// Check result repository methods
	CreateCheckResult(result models.CheckResult) error
	GetCheckResultsByEndpointID(endpointID string, start, end time.Time) ([]models.CheckResult, error)
	GetLatestCheckResultByEndpointID(endpointID string) (models.CheckResult, error)

	// Alert repository methods
	CreateAlert(alert models.Alert) error
	GetAlertsByEndpointID(endpointID string) ([]models.Alert, error)
	GetAlertsByUserID(userID string) ([]models.Alert, error)
	UpdateAlert(alert models.Alert) error
}

// MemoryRepository is an in-memory implementation of Repository
// This is a simple implementation for development purposes
// In production, this would be replaced with a database implementation
type MemoryRepository struct {
	users        map[string]models.User
	usersByEmail map[string]string // email -> id
	endpoints    map[string]models.Endpoint
	checkResults map[string][]models.CheckResult
	alerts       map[string]models.Alert
}

// NewMemoryRepository creates a new in-memory repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users:        make(map[string]models.User),
		usersByEmail: make(map[string]string),
		endpoints:    make(map[string]models.Endpoint),
		checkResults: make(map[string][]models.CheckResult),
		alerts:       make(map[string]models.Alert),
	}
}

// User repository methods

// CreateUser adds a new user to the repository
func (r *MemoryRepository) CreateUser(user models.User) error {
	r.users[user.ID] = user
	r.usersByEmail[user.Email] = user.ID
	return nil
}

// GetUserByID retrieves a user by ID
func (r *MemoryRepository) GetUserByID(id string) (models.User, error) {
	user, exists := r.users[id]
	if !exists {
		return models.User{}, ErrNotFound
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (r *MemoryRepository) GetUserByEmail(email string) (models.User, error) {
	id, exists := r.usersByEmail[email]
	if !exists {
		return models.User{}, ErrNotFound
	}
	return r.GetUserByID(id)
}

// UpdateUser updates an existing user
func (r *MemoryRepository) UpdateUser(user models.User) error {
	_, exists := r.users[user.ID]
	if !exists {
		return ErrNotFound
	}

	// Update email mapping if email has changed
	oldUser := r.users[user.ID]
	if oldUser.Email != user.Email {
		delete(r.usersByEmail, oldUser.Email)
		r.usersByEmail[user.Email] = user.ID
	}

	r.users[user.ID] = user
	return nil
}

// DeleteUser removes a user from the repository
func (r *MemoryRepository) DeleteUser(id string) error {
	user, exists := r.users[id]
	if !exists {
		return ErrNotFound
	}

	delete(r.usersByEmail, user.Email)
	delete(r.users, id)
	return nil
}

// Endpoint repository methods

// CreateEndpoint adds a new endpoint to the repository
func (r *MemoryRepository) CreateEndpoint(endpoint models.Endpoint) error {
	r.endpoints[endpoint.ID] = endpoint
	return nil
}

// GetEndpointByID retrieves an endpoint by ID
func (r *MemoryRepository) GetEndpointByID(id string) (models.Endpoint, error) {
	endpoint, exists := r.endpoints[id]
	if !exists {
		return models.Endpoint{}, ErrNotFound
	}
	return endpoint, nil
}

// GetEndpointsByUserID retrieves all endpoints for a user
func (r *MemoryRepository) GetEndpointsByUserID(userID string) ([]models.Endpoint, error) {
	var results []models.Endpoint

	for _, endpoint := range r.endpoints {
		if endpoint.UserID == userID {
			results = append(results, endpoint)
		}
	}

	return results, nil
}

// UpdateEndpoint updates an existing endpoint
func (r *MemoryRepository) UpdateEndpoint(endpoint models.Endpoint) error {
	_, exists := r.endpoints[endpoint.ID]
	if !exists {
		return ErrNotFound
	}

	r.endpoints[endpoint.ID] = endpoint
	return nil
}

// DeleteEndpoint removes an endpoint from the repository
func (r *MemoryRepository) DeleteEndpoint(id string) error {
	_, exists := r.endpoints[id]
	if !exists {
		return ErrNotFound
	}

	delete(r.endpoints, id)
	return nil
}

// Check result repository methods

// CreateCheckResult adds a new check result to the repository
func (r *MemoryRepository) CreateCheckResult(result models.CheckResult) error {
	r.checkResults[result.EndpointID] = append(r.checkResults[result.EndpointID], result)
	return nil
}

// GetCheckResultsByEndpointID retrieves check results for an endpoint in a date range
func (r *MemoryRepository) GetCheckResultsByEndpointID(endpointID string, start, end time.Time) ([]models.CheckResult, error) {
	results := r.checkResults[endpointID]
	var filtered []models.CheckResult

	for _, result := range results {
		if result.Timestamp.After(start) && result.Timestamp.Before(end) {
			filtered = append(filtered, result)
		}
	}

	return filtered, nil
}

// GetLatestCheckResultByEndpointID retrieves the most recent check result for an endpoint
func (r *MemoryRepository) GetLatestCheckResultByEndpointID(endpointID string) (models.CheckResult, error) {
	results := r.checkResults[endpointID]
	if len(results) == 0 {
		return models.CheckResult{}, ErrNotFound
	}

	// Find the most recent result
	latestResult := results[0]
	for _, result := range results {
		if result.Timestamp.After(latestResult.Timestamp) {
			latestResult = result
		}
	}

	return latestResult, nil
}

// Alert repository methods

// CreateAlert adds a new alert to the repository
func (r *MemoryRepository) CreateAlert(alert models.Alert) error {
	r.alerts[alert.ID] = alert
	return nil
}

// GetAlertsByEndpointID retrieves alerts for an endpoint
func (r *MemoryRepository) GetAlertsByEndpointID(endpointID string) ([]models.Alert, error) {
	var results []models.Alert

	for _, alert := range r.alerts {
		if alert.EndpointID == endpointID {
			results = append(results, alert)
		}
	}

	return results, nil
}

// GetAlertsByUserID retrieves alerts for a user
func (r *MemoryRepository) GetAlertsByUserID(userID string) ([]models.Alert, error) {
	var results []models.Alert

	for _, alert := range r.alerts {
		if alert.UserID == userID {
			results = append(results, alert)
		}
	}

	return results, nil
}

// UpdateAlert updates an existing alert
func (r *MemoryRepository) UpdateAlert(alert models.Alert) error {
	_, exists := r.alerts[alert.ID]
	if !exists {
		return ErrNotFound
	}

	r.alerts[alert.ID] = alert
	return nil
}
