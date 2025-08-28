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

func (l *Link) GetID() int {
	return l.ID
}

func (l *Link) GetSource() int {
	return l.Source
}

func (l *Link) GetDestination() int {
	return l.Destination
}

func (l *Link) GetSlotsByBand(band int) []bool {
	return l.Capacities.Bands[band].Slots
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
	return nil
}

func (l *Link) ReleaseConnection(initial_slot int, size int, band int) error {
	capacity := l.GetSlotsByBand(band)
	if initial_slot+size > len(capacity) {
		slog.Error("Out of range")
		return errors.New("out of range for releasing connection")
	}
	for i := initial_slot; i < initial_slot+size; i++ {
		capacity[i] = false
	}
	return nil
}
