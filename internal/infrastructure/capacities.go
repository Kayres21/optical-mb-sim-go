package infrastructure

import (
	"encoding/json"
	"log"
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
		log.Fatalf("Error opening file: %v", err)
		return Capacity{}, err
	}

	var capacities Capacity

	err = json.Unmarshal(dataBytes, &capacities)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return Capacity{}, err
	}

	for i := range capacities.Bands {

		slots := make([]bool, capacities.Bands[i].SlotsLen)

		capacities.Bands[i].Slots = slots
	}
	return capacities, nil
}
