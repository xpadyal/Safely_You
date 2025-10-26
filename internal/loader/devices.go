// Package loader handles loading device data from external sources.
// Currently supports loading device IDs from CSV files for initialization.
package loader

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/xpadyal/Safely_You/internal/models"
	"github.com/xpadyal/Safely_You/internal/store"
)

// LoadDevicesFromCSV reads device IDs from a CSV file and initializes them in the store
func LoadDevicesFromCSV(filename string, storeInstance *models.Store) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header, process device rows
	for _, record := range records[1:] {
		if len(record) > 0 && record[0] != "" {
			store.EnsureDevice(storeInstance, record[0])
		}
	}

	fmt.Printf("Loaded devices from %s\n", filename)
	return nil
}
