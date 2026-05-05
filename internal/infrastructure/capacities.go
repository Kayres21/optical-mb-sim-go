package infrastructure

import (
	"encoding/json"
	"fmt"
	"os"
)

type Capacity struct {
	Bands []Band `json:"bands"`
}

type Band struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SlotsLen int    `json:"capacity"`
	Slots    []bool `json:"-"`
}

func ReadCapacityFile(capacityPath string) (Capacity, error) {
	dataBytes, err := os.ReadFile(capacityPath)
	if err != nil {
		return Capacity{}, fmt.Errorf("reading capacity file %q: %w", capacityPath, err)
	}

	var capacities Capacity
	if err = json.Unmarshal(dataBytes, &capacities); err != nil {
		return Capacity{}, fmt.Errorf("parsing capacity file %q: %w", capacityPath, err)
	}

	for i := range capacities.Bands {
		capacities.Bands[i].Slots = make([]bool, capacities.Bands[i].SlotsLen)
	}

	return capacities, nil
}
