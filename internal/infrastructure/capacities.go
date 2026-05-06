package infrastructure

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/Kayres21/optical-mb-sim-go/pkg/validator"
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
	schemaPath := filepath.Join(filepath.Dir(capacityPath), "schema.json")
	dataBytes, err := validator.ValidateFile(capacityPath, schemaPath)
	if err != nil {
		return Capacity{}, fmt.Errorf("validating capacity file: %w", err)
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
