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

func (l *Link) AssignConnection(initialSlot, slotCount, band int) error {
	capacity := l.GetSlotsByBand(band)
	if initialSlot+slotCount > len(capacity) {
		slog.Error("AssignConnection: out of range", "initialSlot", initialSlot, "slotCount", slotCount)
		return errors.New("out of range for assigning connection")
	}
	for i := initialSlot; i < initialSlot+slotCount; i++ {
		capacity[i] = true
	}
	return nil
}

func (l *Link) ReleaseConnection(initialSlot, slotCount, band int) error {
	capacity := l.GetSlotsByBand(band)
	if initialSlot+slotCount > len(capacity) {
		slog.Error("ReleaseConnection: out of range", "initialSlot", initialSlot, "slotCount", slotCount)
		return errors.New("out of range for releasing connection")
	}
	for i := initialSlot; i < initialSlot+slotCount; i++ {
		capacity[i] = false
	}
	return nil
}
