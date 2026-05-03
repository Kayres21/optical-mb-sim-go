package connections

import "testing"

func TestConnectionFields(t *testing.T) {
	conn := Connection{}

	conn.Id = "conn1"
	if conn.Id != "conn1" {
		t.Errorf("Expected Id to be 'conn1', got '%s'", conn.Id)
	}

	conn.Source = 1
	if conn.Source != 1 {
		t.Errorf("Expected Source to be 1, got %d", conn.Source)
	}

	conn.Destination = 2
	if conn.Destination != 2 {
		t.Errorf("Expected Destination to be 2, got %d", conn.Destination)
	}

	conn.Slots = 10
	if conn.Slots != 10 {
		t.Errorf("Expected Slots to be 10, got %d", conn.Slots)
	}

	conn.InitialSlot = 5
	if conn.InitialSlot != 5 {
		t.Errorf("Expected InitialSlot to be 5, got %d", conn.InitialSlot)
	}

	conn.FinalSlot = 15
	if conn.FinalSlot != 15 {
		t.Errorf("Expected FinalSlot to be 15, got %d", conn.FinalSlot)
	}

	conn.BandSelected = 3
	if conn.BandSelected != 3 {
		t.Errorf("Expected BandSelected to be 3, got %d", conn.BandSelected)
	}

	conn.Allocated = true
	if !conn.Allocated {
		t.Errorf("Expected Allocated to be true, got false")
	}
}
