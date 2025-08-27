package infrastructure

import (
	"errors"
	"log/slog"
)

type Link struct {
	ID          int      `json:"id"`
	Source      int      `json:"src"`
	Destination int      `json:"dst"`
	Length      int      `json:"length"`
	Capacities  Capacity `json:"-"`
}

func (l *Link) GetSlotsByBand(band int) []bool {
	return l.Capacities.Bands[band].Slots
}

func (l *Link) AssignCapacityByBand(capacity []bool, band int) {
	l.Capacities.Bands[band].Slots = capacity
}

func (l *Link) AssignConnection(initial_slot int, size int, band int) error {

	capacity := l.GetSlotsByBand(band)

	if initial_slot+size > len(capacity) {
		slog.Error("Out of range")
		return errors.New("out of range for assigning connection")
	}

	for i := initial_slot; i < initial_slot+size; i++ {
		capacity[i] = true
	}

	l.AssignCapacityByBand(capacity, band)

	return nil

}
func (l *Link) ReleaseConnection(initial_slot int, size int, band int) error {
	capacity := l.GetSlotsByBand(band)

	if initial_slot+size > len(capacity) {
		slog.Error("Out of range")
		return errors.New("out of range for assigning connection")
	}

	for i := initial_slot; i < initial_slot+size; i++ {
		capacity[i] = false
	}

	l.AssignCapacityByBand(capacity, band)
	return nil
}
