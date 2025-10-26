package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xpadyal/Safely_You/internal/models"
	"github.com/xpadyal/Safely_You/internal/store"
	"github.com/xpadyal/Safely_You/internal/utils"
)

// Response helper functions for common error responses
func notFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, models.NotFoundResponse{Msg: msg})
}

func badRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, models.ErrorResponse{Msg: msg})
}

func internalError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, models.ErrorResponse{Msg: msg})
}

// Handles heartbeat registration from devices
func PostHeartbeatHandler(storeInstance *models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if deviceID == "" {
			notFound(c, "device not found")
			return
		}

		var req models.HeartbeatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			badRequest(c, "invalid JSON body")
			return
		}

		t, err := utils.ParseRFC3339(req.SentAt)
		if err != nil {
			badRequest(c, "invalid sent_at: "+err.Error())
			return
		}

		// Check if device exists, return 404 if not
		if _, exists := store.SnapshotDevice(storeInstance, deviceID); !exists {
			notFound(c, "device not found")
			return
		}

		// Add heartbeat with error handling
		if err := store.AddHeartbeat(storeInstance, deviceID, t); err != nil {
			internalError(c, "failed to add heartbeat")
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// Handles upload stats submission
func PostStatsHandler(storeInstance *models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if deviceID == "" {
			notFound(c, "device not found")
			return
		}

		var req models.UploadStatsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			badRequest(c, "invalid JSON body")
			return
		}

		// Check sent_at format even though we don't use it for calculations
		if _, err := utils.ParseRFC3339(req.SentAt); err != nil {
			badRequest(c, "invalid sent_at: "+err.Error())
			return
		}

		// Check if device exists, return 404 if not
		if _, exists := store.SnapshotDevice(storeInstance, deviceID); !exists {
			notFound(c, "device not found")
			return
		}

		// Add upload time with error handling
		if err := store.AddUploadTime(storeInstance, deviceID, req.UploadTime); err != nil {
			internalError(c, "failed to add upload time")
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// Returns device statistics
func GetStatsHandler(storeInstance *models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if deviceID == "" {
			notFound(c, "device not found")
			return
		}

		d, exists := store.SnapshotDevice(storeInstance, deviceID)
		if !exists {
			notFound(c, "device not found")
			return
		}

		// Check if device has any data, return 204 if no data
		if len(d.Heartbeats) == 0 && len(d.UploadTimes) == 0 {
			c.Status(http.StatusNoContent)
			return
		}

		// Compute stats with error handling
		uptime, err := store.ComputeUptime(d)
		if err != nil {
			internalError(c, "failed to compute uptime")
			return
		}

		avgUpload, err := store.ComputeAvgUpload(d)
		if err != nil {
			internalError(c, "failed to compute average upload time")
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
