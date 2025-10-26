// Package handlers contains HTTP request handlers for the device monitoring API.
// It provides endpoints for heartbeat registration, upload stats submission, and device statistics retrieval.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xpadyal/Safely_You/internal/models"
	"github.com/xpadyal/Safely_You/internal/store"
	"github.com/xpadyal/Safely_You/internal/utils"
	"github.com/xpadyal/Safely_You/internal/validation"
)

// Handles heartbeat registration from devices
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

// Handles upload stats submission
func PostStatsHandler(storeInstance *models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if !validation.ValidateDeviceExists(c, storeInstance, deviceID) {
			return
		}

		var req models.UploadStatsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			validation.BadRequest(c, "invalid JSON body")
			return
		}

		// Check sent_at format and validate reasonableness
		_, ok := validation.ValidateAndExtractTimestamp(c, req.SentAt)
		if !ok {
			return
		}

		// Add upload time with error handling
		if err := store.AddUploadTime(storeInstance, deviceID, req.UploadTime); err != nil {
			validation.InternalError(c, "failed to add upload time")
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// Returns device statistics
func GetStatsHandler(storeInstance *models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if !validation.ValidateDeviceExists(c, storeInstance, deviceID) {
			return
		}

		d, _ := store.SnapshotDevice(storeInstance, deviceID) // We know it exists from validation

		// Check if device has any data, return 204 if no data
		if len(d.Heartbeats) == 0 && len(d.UploadTimes) == 0 {
			c.Status(http.StatusNoContent)
			return
		}

		// Compute stats with error handling
		uptime, err := store.ComputeUptime(d)
		if err != nil {
			validation.InternalError(c, "failed to compute uptime")
			return
		}

		avgUpload, err := store.ComputeAvgUpload(d)
		if err != nil {
			validation.InternalError(c, "failed to compute average upload time")
			return
		}

		resp := models.StatsGetResponse{
			Uptime:        utils.Round2(uptime),
			AvgUploadTime: avgUpload,
		}

		c.JSON(http.StatusOK, resp)
	}
}

// Simple health check endpoint
func HealthHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", []byte("OK"))
}
