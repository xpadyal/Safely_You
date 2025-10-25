package store

import (
	"errors"
	"fmt"
	"time"

	"github.com/xpadyal/Safely_You/internal/models"
)

type Store struct {
	devices map[string]*models.Device
}

func NewStore() *Store {
	return &Store{
		devices: make(map[string]*models.Device),
	}
}

// EnsureDevice creates a device if it doesn't exist and returns it
func (s *Store) EnsureDevice(id string) *models.Device {
	if d, ok := s.devices[id]; ok {
		return d
	}
	d := &models.Device{}
	s.devices[id] = d
	return d
}

// Helper functions

func ParseRFC3339(ts string) (time.Time, error) {
	// Parse RFC3339 format as specified in OpenAPI contract
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid sent_at format: %s", ts)
	}
	return t.UTC(), nil
}

// Counts unique minute buckets in UTC
// Example: 10:00:05 and 10:00:55 are in the same minute
func UniqueMinuteCount(times []time.Time) int {
	if len(times) == 0 {
		return 0
	}
	seen := make(map[string]struct{}, len(times))
	for _, t := range times {
		key := t.UTC().Format("2006-01-02T15:04")
		seen[key] = struct{}{}
	}
	return len(seen)
}

// Calculates the inclusive minute span between first and last timestamps
// This includes both the first and last minute buckets
func MinutesBetweenFirstAndLast(times []time.Time) (int, error) {
	if len(times) < 2 {
		return 0, errors.New("need at least two timestamps to form a window")
	}
	minT := times[0]
	maxT := times[0]
	for _, t := range times[1:] {
		if t.Before(minT) {
			minT = t
		}
		if t.After(maxT) {
			maxT = t
		}
	}
	if !maxT.After(minT) {
		return 0, errors.New("non-positive window")
	}
	span := int(maxT.Sub(minT)/time.Minute) + 1
	return span, nil
}

// Formula: (unique_minute_heartbeats / total_minute_span) * 100
func ComputeUptime(d *models.Device) float64 {
	if d == nil || len(d.Heartbeats) == 0 {
		return 0.0
	}
	uniq := UniqueMinuteCount(d.Heartbeats)
	windowMinutes, err := MinutesBetweenFirstAndLast(d.Heartbeats)
	if err != nil || windowMinutes == 0 {
		// All heartbeats in the same minute = 100% uptime
		if uniq > 0 {
			return 100.0
		}
		return 0.0
	}
	return (float64(uniq) / float64(windowMinutes)) * 100.0
}

// Calculates average upload time and returns it as a duration string
// Returns ("0s", false) if no data available
func ComputeAvgUpload(d *models.Device) (string, bool) {
	if d == nil || len(d.UploadTimes) == 0 {
		return "0s", false
	}
	var sum int64
	for _, v := range d.UploadTimes {
		sum += v
	}
	avg := sum / int64(len(d.UploadTimes))
	// Convert nanoseconds to duration and format as string
	duration := time.Duration(avg)
	return duration.String(), true
}

// Device operations

func (s *Store) AddHeartbeat(id string, t time.Time) {
	d, ok := s.devices[id]
	if !ok {
		d = &models.Device{}
		s.devices[id] = d
	}
	d.Heartbeats = append(d.Heartbeats, t.UTC())
}

func (s *Store) AddUploadTime(id string, v int64) {
	d, ok := s.devices[id]
	if !ok {
		d = &models.Device{}
		s.devices[id] = d
	}
	d.UploadTimes = append(d.UploadTimes, v)
}

// Returns the device data
func (s *Store) SnapshotDevice(id string) (*models.Device, bool) {
	d, ok := s.devices[id]
	return d, ok
}
