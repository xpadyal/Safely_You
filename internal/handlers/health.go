// Package handlers contains HTTP request handlers for the device monitoring API.
// It provides endpoints for heartbeat registration, upload stats submission, and device statistics retrieval.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler provides a simple health check endpoint.
// It returns a 200 OK status with "OK" in the response body.
func HealthHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", []byte("OK"))
}
