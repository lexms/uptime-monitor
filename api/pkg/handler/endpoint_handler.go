package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tou01/uptime-monitor/pkg/middleware"
	"github.com/tou01/uptime-monitor/pkg/models"
	"github.com/tou01/uptime-monitor/pkg/service"
)

// EndpointHandler handles HTTP requests related to monitoring endpoints
type EndpointHandler struct {
	monitorService *service.MonitorService
}

// NewEndpointHandler creates a new endpoint handler
func NewEndpointHandler(monitorService *service.MonitorService) *EndpointHandler {
	return &EndpointHandler{
		monitorService: monitorService,
	}
}

// RegisterRoutes registers the endpoint routes to the router
func (h *EndpointHandler) RegisterRoutes(router *gin.RouterGroup) {
	endpoints := router.Group("/endpoints")
	{
		endpoints.GET("", h.GetAllEndpoints)
		endpoints.GET("/:id", h.GetEndpoint)
		endpoints.POST("", h.CreateEndpoint)
		endpoints.PUT("/:id", h.UpdateEndpoint)
		endpoints.DELETE("/:id", h.DeleteEndpoint)
	}
}

// GetAllEndpoints retrieves all endpoints for the current user
func (h *EndpointHandler) GetAllEndpoints(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)

	// Get all endpoints for the user
	endpoints := h.monitorService.GetAllEndpoints(userID)

	c.JSON(http.StatusOK, gin.H{
		"data": endpoints,
	})
}

// GetEndpoint retrieves a single endpoint by ID
func (h *EndpointHandler) GetEndpoint(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	// Get the endpoint
	endpoint, err := h.monitorService.GetEndpoint(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if endpoint belongs to the user
	if endpoint.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this endpoint",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": endpoint,
	})
}

// CreateEndpoint creates a new endpoint to monitor
func (h *EndpointHandler) CreateEndpoint(c *gin.Context) {
	var endpoint models.Endpoint

	// Bind JSON to endpoint
	if err := c.ShouldBindJSON(&endpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set user ID
	endpoint.UserID = middleware.GetUserID(c)

	// Create endpoint
	if err := h.monitorService.AddEndpoint(endpoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": endpoint,
	})
}

// UpdateEndpoint updates an existing endpoint
func (h *EndpointHandler) UpdateEndpoint(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	// Get the current endpoint
	currentEndpoint, err := h.monitorService.GetEndpoint(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if endpoint belongs to the user
	if currentEndpoint.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this endpoint",
		})
		return
	}

	// Bind JSON to endpoint
	var updatedEndpoint models.Endpoint
	if err := c.ShouldBindJSON(&updatedEndpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set ID and user ID
	updatedEndpoint.ID = id
	updatedEndpoint.UserID = userID

	// Update endpoint
	if err := h.monitorService.UpdateEndpoint(updatedEndpoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": updatedEndpoint,
	})
}

// DeleteEndpoint removes an endpoint from monitoring
func (h *EndpointHandler) DeleteEndpoint(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	// Get the current endpoint
	currentEndpoint, err := h.monitorService.GetEndpoint(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if endpoint belongs to the user
	if currentEndpoint.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this endpoint",
		})
		return
	}

	// Delete endpoint
	if err := h.monitorService.DeleteEndpoint(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Endpoint deleted successfully",
	})
}
