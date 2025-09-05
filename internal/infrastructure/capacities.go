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

func (c *Capacity) GetBands() []Band {
	return c.Bands
}

func (c *Capacity) SetBands(bands []Band) {
	c.Bands = bands
}

func (b *Band) GetID() string {
	return b.ID
}

func (b *Band) GetSlotsLen() int {
	return b.SlotsLen
}

func (b *Band) GetName() string {
	return b.Name
}

func (b *Band) GetSlots() []bool {
	return b.Slots
}

func (b *Band) SetSlots(slots []bool) {
	b.Slots = slots
}

func (b *Band) SetSlotsLen(slotsLen int) {
	b.SlotsLen = slotsLen
}

func (b *Band) SetName(name string) {
	b.Name = name
}

func (b *Band) SetID(id string) {
	b.ID = id
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

		slots := make([]bool, capacities.Bands[i].GetSlotsLen())

		capacities.Bands[i].SetSlots(slots)
	}
	return capacities, nil
}
