// Package validation provides validation helpers for HTTP handlers.
// It contains common validation logic used across multiple endpoints.
package validation

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xpadyal/Safely_You/internal/models"
	"github.com/xpadyal/Safely_You/internal/store"
	"github.com/xpadyal/Safely_You/internal/utils"
)

// Response helper functions for common error responses
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, models.NotFoundResponse{Msg: msg})
}

func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, models.ErrorResponse{Msg: msg})
}

func InternalError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, models.ErrorResponse{Msg: msg})
}

// ValidateAndExtractTimestamp validates and parses a timestamp string
func ValidateAndExtractTimestamp(c *gin.Context, sentAt string) (time.Time, bool) {
	t, err := utils.ParseRFC3339(sentAt)
	if err != nil {
		BadRequest(c, "invalid sent_at: "+err.Error())
		return time.Time{}, false
	}

	if err := utils.ValidateTimestamp(t); err != nil {
		BadRequest(c, err.Error())
		return time.Time{}, false
	}

	return t, true
}

// ValidateDeviceExists checks if a device exists and validates the device ID
func ValidateDeviceExists(c *gin.Context, storeInstance *models.Store, deviceID string) bool {
	if deviceID == "" {
		NotFound(c, "device not found")
		return false
	}

	if _, exists := store.SnapshotDevice(storeInstance, deviceID); !exists {
		NotFound(c, "device not found")
		return false
	}

	return true
}
