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
	// FragmentationRatioByBand stores FR1 per band and is updated on link usage.
	FragmentationRatioByBand []float64 `json:"-"`
	// FragmentationRatio2ByBand stores FR2 per band and is updated on link usage.
	FragmentationRatio2ByBand []float64 `json:"-"`
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
	l.UpdateFragmentationRatio(band)
	l.UpdateFragmentationRatio2(band)
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
	l.UpdateFragmentationRatio(band)
	l.UpdateFragmentationRatio2(band)
	return nil
}

func (l *Link) GetFragmentationRatioByBand(band int) float64 {
	if band < 0 || band >= len(l.FragmentationRatioByBand) {
		return 0
	}
	return l.FragmentationRatioByBand[band]
}

func (l *Link) GetFragmentationRatio2ByBand(band int) float64 {
	if band < 0 || band >= len(l.FragmentationRatio2ByBand) {
		return 0
	}
	return l.FragmentationRatio2ByBand[band]
}

func (l *Link) UpdateAllFragmentationRatios() {
	if len(l.Capacities.Bands) == 0 {
		l.FragmentationRatioByBand = nil
		l.FragmentationRatio2ByBand = nil
		return
	}

	if len(l.FragmentationRatioByBand) != len(l.Capacities.Bands) {
		l.FragmentationRatioByBand = make([]float64, len(l.Capacities.Bands))
	}
	if len(l.FragmentationRatio2ByBand) != len(l.Capacities.Bands) {
		l.FragmentationRatio2ByBand = make([]float64, len(l.Capacities.Bands))
	}

	for band := range l.Capacities.Bands {
		l.UpdateFragmentationRatio(band)
		l.UpdateFragmentationRatio2(band)
	}
}

func (l *Link) UpdateFragmentationRatio(band int) {
	if band < 0 || band >= len(l.Capacities.Bands) {
		return
	}
	if len(l.FragmentationRatioByBand) != len(l.Capacities.Bands) {
		l.FragmentationRatioByBand = make([]float64, len(l.Capacities.Bands))
	}
	if len(l.FragmentationRatio2ByBand) != len(l.Capacities.Bands) {
		l.FragmentationRatio2ByBand = make([]float64, len(l.Capacities.Bands))
	}

	slots := l.GetSlotsByBand(band)
	if len(slots) == 0 {
		l.FragmentationRatioByBand[band] = 0
		return
	}

	maxBlock := 0
	currentBlock := 0
	availableSlots := 0

	for _, occupied := range slots {
		if occupied {
			currentBlock++
			if currentBlock > maxBlock {
				maxBlock = currentBlock
			}
			continue
		}

		availableSlots++
		currentBlock = 0
	}

	totalSlots := len(slots)
	usedSlots := totalSlots - availableSlots
	if usedSlots <= 0 {
		l.FragmentationRatioByBand[band] = 0
		return
	}

	l.FragmentationRatioByBand[band] = 1 - (float64(maxBlock) / float64(usedSlots))
}

func (l *Link) UpdateFragmentationRatio2(band int) {
	if band < 0 || band >= len(l.Capacities.Bands) {
		return
	}
	if len(l.FragmentationRatio2ByBand) != len(l.Capacities.Bands) {
		l.FragmentationRatio2ByBand = make([]float64, len(l.Capacities.Bands))
	}

	slots := l.GetSlotsByBand(band)
	if len(slots) == 0 {
		l.FragmentationRatio2ByBand[band] = 0
		return
	}

	freeBlocks := 0
	inFreeBlock := false
	for _, occupied := range slots {
		if !occupied {
			if !inFreeBlock {
				freeBlocks++
				inFreeBlock = true
			}
			continue
		}
		inFreeBlock = false
	}

	if freeBlocks <= 0 {
		l.FragmentationRatio2ByBand[band] = 0
		return
	}

	l.FragmentationRatio2ByBand[band] = 1 - (1 / float64(freeBlocks))
}
