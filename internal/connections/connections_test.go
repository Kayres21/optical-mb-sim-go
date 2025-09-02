package connections

import "testing"

func TestConnectionGettersAndSetters(t *testing.T) {
	conn := Connection{}

	conn.SetId("conn1")
	if conn.GetId() != "conn1" {
		t.Errorf("Expected Id to be 'conn1', got '%s'", conn.GetId())
	}

	conn.SetSource(1)
	if conn.GetSource() != 1 {
		t.Errorf("Expected Source to be 1, got %d", conn.GetSource())
	}

	conn.SetDestination(2)
	if conn.GetDestination() != 2 {
		t.Errorf("Expected Destination to be 2, got %d", conn.GetDestination())
	}

	conn.SetSlots(10)
	if conn.GetSlots() != 10 {
		t.Errorf("Expected Slots to be 10, got %d", conn.GetSlots())
	}

	conn.SetInitialSlot(5)
	if conn.GetInitialSlot() != 5 {
		t.Errorf("Expected InitialSlot to be 5, got %d", conn.GetInitialSlot())
	}

	conn.SetFinalSlot(15)
	if conn.GetFinalSlot() != 15 {
		t.Errorf("Expected FinalSlot to be 15, got %d", conn.GetFinalSlot())
	}

	conn.SetBandSelected(3)
	if conn.GetBandSelected() != 3 {
		t.Errorf("Expected BandSelected to be 3, got %d", conn.GetBandSelected())
	}

	conn.Allocated = true
	if !conn.IsAllocated() {
		t.Errorf("Expected Allocated to be true, got false")
	}
}
