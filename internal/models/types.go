// Package models defines the core data structures used throughout the application.
// It includes device data models, request/response types, and the in-memory store.
package models

import (
	"sync"
	"time"
)

type Device struct {
	Heartbeats  []time.Time
	UploadTimes []int64 // upload times in nanoseconds
}

type HeartbeatRequest struct {
	SentAt string `json:"sent_at"`
}

type UploadStatsRequest struct {
	SentAt     string `json:"sent_at"`
	UploadTime int64  `json:"upload_time"` // upload time in nanoseconds
}

type StatsGetResponse struct {
	Uptime        float64 `json:"uptime"`
	AvgUploadTime string  `json:"avg_upload_time"`
}

// Response types for error handling
type NotFoundResponse struct {
	Msg string `json:"msg"`
}

type ErrorResponse struct {
	Msg string `json:"msg"`
}

// Store holds device data in memory with thread-safe access
type Store struct {
	Mu      sync.RWMutex
	Devices map[string]*Device
}
