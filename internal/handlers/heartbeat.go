// Package handlers contains HTTP request handlers for the device monitoring API.
// It provides endpoints for heartbeat registration, upload stats submission, and device statistics retrieval.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xpadyal/Safely_You/internal/models"
	"github.com/xpadyal/Safely_You/internal/store"
	"github.com/xpadyal/Safely_You/internal/validation"
)

// PostHeartbeatHandler handles heartbeat registration from devices.
// It validates the device ID and request payload, then adds the heartbeat to the store.
func PostHeartbeatHandler(storeInstance *models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if !validation.ValidateDeviceExists(c, storeInstance, deviceID) {
			return
		}

		var req models.HeartbeatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			validation.BadRequest(c, "invalid JSON body")
			return
		}

		t, ok := validation.ValidateAndExtractTimestamp(c, req.SentAt)
		if !ok {
			return
		}

		// Add heartbeat with error handling
		if err := store.AddHeartbeat(storeInstance, deviceID, t); err != nil {
			validation.InternalError(c, "failed to add heartbeat")
			return
		}

		c.Status(http.StatusNoContent)
	}
}
