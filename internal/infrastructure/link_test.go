package infrastructure

import "testing"

func TestLinkGetID(t *testing.T) {
	link := Link{ID: 1}
	if link.GetID() != 1 {
		t.Errorf("Expected ID to be 1, got %d", link.GetID())
	}
}

func TestLinkGetSource(t *testing.T) {
	link := Link{Source: 2}
	if link.GetSource() != 2 {
		t.Errorf("Expected Source to be 2, got %d", link.GetSource())
	}
}

func TestLinkGetDestination(t *testing.T) {
	link := Link{Destination: 3}
	if link.GetDestination() != 3 {
		t.Errorf("Expected Destination to be 3, got %d", link.GetDestination())
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
}
