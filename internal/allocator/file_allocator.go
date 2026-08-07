package allocator

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

type fileAllocationEvent struct {
	ID          string
	Source      int
	Destination int
	Allocated   bool
}

// FirstFitFromFile loads a CSV file with allocation decisions and returns an allocator
// function that matches incoming requests by ID, source, and destination.
func FirstFitFromFile(csvPath string) (Allocator, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("open first fit CSV: %w", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read first fit CSV: %w", err)
	}

	events := make([]fileAllocationEvent, 0, len(records)-1)
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 4 {
			return nil, fmt.Errorf("invalid first fit CSV row %d: expected 4 columns", i+1)
		}

		id := strings.TrimSpace(record[0])
		source, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid source value on row %d: %w", i+1, err)
		}
		destination, err := strconv.Atoi(strings.TrimSpace(record[2]))
		if err != nil {
			return nil, fmt.Errorf("invalid destination value on row %d: %w", i+1, err)
		}

		allocated := strings.EqualFold(strings.TrimSpace(record[3]), "ALLOCATED")

		events = append(events, fileAllocationEvent{
			ID:          id,
			Source:      source,
			Destination: destination,
			Allocated:   allocated,
		})
	}

	return func(source, destination int, getSlot func(band int) int, network infrastructure.Network, path connections.Routes, numberOfBands int, id string, addConnection func(connections.Connection)) bool {
		for index, event := range events {
			if event.ID == id && event.Source == source && event.Destination == destination {
				events = append(events[:index], events[index+1:]...)
				if !event.Allocated {
					return false
				}

				return true
			}
		}
		return false
	}, nil
}
