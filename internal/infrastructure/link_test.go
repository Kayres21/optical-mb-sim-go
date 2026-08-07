package infrastructure

import (
	"math"
	"testing"
)

func TestLinkGetID(t *testing.T) {
	link := Link{ID: 1}
	if link.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", link.ID)
	}
}

func TestLinkGetSource(t *testing.T) {
	link := Link{Source: 2}
	if link.Source != 2 {
		t.Errorf("Expected Source to be 2, got %d", link.Source)
	}
}

func TestLinkGetDestination(t *testing.T) {
	link := Link{Destination: 3}
	if link.Destination != 3 {
		t.Errorf("Expected Destination to be 3, got %d", link.Destination)
	}
}

func TestLinkGetSlotsByBand(t *testing.T) {
	link := Link{
		Capacities: Capacity{
			Bands: []Band{
				{Slots: []bool{false, true, false}},
				{Slots: []bool{true, true, true}},
			},
		},
	}
	slots := link.GetSlotsByBand(1)
	expected := []bool{true, true, true}
	for i, slot := range slots {
		if slot != expected[i] {
			t.Errorf("Expected slot %d to be %v, got %v", i, expected[i], slot)
		}
	}
}

func TestLinkAssignAndReleaseConnection(t *testing.T) {
	link := Link{
		Capacities: Capacity{
			Bands: []Band{
				{Slots: make([]bool, 10)},
			},
		},
	}

	err := link.AssignConnection(2, 3, 0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	for i := 2; i < 5; i++ {
		if !link.GetSlotsByBand(0)[i] {
			t.Errorf("Expected slot %d to be true, got false", i)
		}
	}

	err = link.ReleaseConnection(2, 3, 0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	for i := 2; i < 5; i++ {
		if link.GetSlotsByBand(0)[i] {
			t.Errorf("Expected slot %d to be false, got true", i)
		}
	}

	if got := link.GetFragmentationRatioByBand(0); got != 0 {
		t.Errorf("Expected fragmentation ratio to be 0 after release, got %f", got)
	}
}

func TestLinkFragmentationRatio1UpdatesByUse(t *testing.T) {
	link := Link{
		Capacities: Capacity{
			Bands: []Band{{Slots: make([]bool, 8)}},
		},
	}

	if err := link.AssignConnection(0, 3, 0); err != nil {
		t.Fatalf("unexpected assign error: %v", err)
	}
	if got := link.GetFragmentationRatioByBand(0); got != 0 {
		t.Fatalf("expected 0 fragmentation with single block, got %f", got)
	}

	if err := link.AssignConnection(5, 1, 0); err != nil {
		t.Fatalf("unexpected assign error: %v", err)
	}

	// FR1: 1 - maxBlock/usedSlots = 1 - 3/4 = 0.25
	expected := 0.25
	if got := link.GetFragmentationRatioByBand(0); math.Abs(got-expected) > 1e-9 {
		t.Fatalf("expected fragmentation ratio %f, got %f", expected, got)
	}

	if err := link.ReleaseConnection(0, 3, 0); err != nil {
		t.Fatalf("unexpected release error: %v", err)
	}
	if got := link.GetFragmentationRatioByBand(0); got != 0 {
		t.Fatalf("expected 0 fragmentation with one used slot, got %f", got)
	}
}

func TestLinkFragmentationRatio2UpdatesByUse(t *testing.T) {
	link := Link{
		Capacities: Capacity{
			Bands: []Band{{Slots: []bool{true, false, false, true, false, false, false, true}}},
		},
	}

	link.UpdateAllFragmentationRatios()
	if got := link.GetFragmentationRatio2ByBand(0); math.Abs(got-0.5) > 1e-9 {
		t.Fatalf("expected FR2 of 0.5 with two free blocks, got %f", got)
	}

	if err := link.ReleaseConnection(0, 1, 0); err != nil {
		t.Fatalf("unexpected release error: %v", err)
	}
	if got := link.GetFragmentationRatio2ByBand(0); math.Abs(got-0.5) > 1e-9 {
		t.Fatalf("expected FR2 to remain 0.5 with two free blocks, got %f", got)
	}

	if err := link.AssignConnection(4, 3, 0); err != nil {
		t.Fatalf("unexpected assign error: %v", err)
	}
	if got := link.GetFragmentationRatio2ByBand(0); got != 0 {
		t.Fatalf("expected FR2 of 0 with a single free block, got %f", got)
	}
}
