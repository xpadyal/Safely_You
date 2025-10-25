package models

import "time"

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
