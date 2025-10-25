package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xpadyal/Safely_You/internal/models"
	"github.com/xpadyal/Safely_You/internal/store"
	"github.com/xpadyal/Safely_You/internal/utils"
)

// Handles heartbeat registration from devices
func PostHeartbeatHandler(storeInstance *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if deviceID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}

		var req models.HeartbeatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
			return
		}

		t, err := store.ParseRFC3339(req.SentAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sent_at: " + err.Error()})
			return
		}

		storeInstance.AddHeartbeat(deviceID, t)
		c.Status(http.StatusNoContent)
	}
}

// Handles upload stats submission
func PostStatsHandler(storeInstance *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if deviceID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}

		var req models.UploadStatsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
			return
		}

		// Check sent_at format even though we don't use it for calculations
		if _, err := store.ParseRFC3339(req.SentAt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sent_at: " + err.Error()})
			return
		}

		storeInstance.AddUploadTime(deviceID, req.UploadTime)
		c.Status(http.StatusNoContent)
	}
}

// Returns device statistics
func GetStatsHandler(storeInstance *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		if deviceID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}

		d, exists := storeInstance.SnapshotDevice(deviceID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}

		uptime := store.ComputeUptime(d)
		avgUpload, _ := store.ComputeAvgUpload(d)
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
