package models

import (
	"time"
)

// Endpoint represents a URL to be monitored for uptime
type Endpoint struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	Method       string    `json:"method"`
	Headers      []Header  `json:"headers"`
	Body         string    `json:"body,omitempty"`
	Interval     int       `json:"interval"` // in seconds
	Timeout      int       `json:"timeout"`  // in seconds
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	UserID       string    `json:"userId"`
	IsActive     bool      `json:"isActive"`
	ExpectStatus int       `json:"expectStatus,omitempty"` // Expected HTTP status code
}

// Header represents an HTTP header
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CheckResult represents the result of a single uptime check
type CheckResult struct {
	ID           string    `json:"id"`
	EndpointID   string    `json:"endpointId"`
	Status       int       `json:"status"`
	StatusText   string    `json:"statusText"`
	ResponseTime int64     `json:"responseTime"` // in milliseconds
	Timestamp    time.Time `json:"timestamp"`
	Success      bool      `json:"success"`
	Error        string    `json:"error,omitempty"`
	ResponseBody string    `json:"responseBody,omitempty"`
}

// Alert represents a notification sent when an endpoint is down
type Alert struct {
	ID         string    `json:"id"`
	EndpointID string    `json:"endpointId"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
	Status     string    `json:"status"` // "triggered", "resolved", "acknowledged"
	UserID     string    `json:"userId"`
}
