package allocator

import (
	"testing"

	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

func TestFirstFit_Success(t *testing.T) {
	// Setup a simple route
	routes := connections.Routes{
		Paths: []connections.Path{
			{
				Source:      0,
				Destination: 1,
				PathLinks:   [][]int{{0, 1}}, // Node IDs: 0, 1
			},
		},
	}

	// Setup a simple network
	// We need a link with ID 10 and enough slots.
	links := []infrastructure.Link{
		{
			ID:          10,
			Source:      0,
			Destination: 1,
			Capacities: infrastructure.Capacity{
				Bands: []infrastructure.Band{
					{Slots: []bool{false, false, false, false}},
				},
			},
		},
	}
	network := infrastructure.Network{Links: links}

	allocated, conn := FirstFit(0, 1, 2, network, routes, 1, "test-id")

	if !allocated {
		t.Fatalf("expected allocation to succeed")
	}

	if conn.Id != "test-id" {
		t.Errorf("expected connection id test-id, got %v", conn.Id)
	}

	if conn.InitialSlot != 0 {
		t.Errorf("expected initial slot 0, got %v", conn.InitialSlot)
	}

	if conn.FinalSlot != 1 {
		t.Errorf("expected final slot 1, got %v", conn.FinalSlot)
	}
}

func TestFirstFit_FailNoCapacity(t *testing.T) {
	routes := connections.Routes{
		Paths: []connections.Path{
			{
				Source:      0,
				Destination: 1,
				PathLinks:   [][]int{{0, 1}},
			},
		},
	}

	// Link is full
	links := []infrastructure.Link{
		{
			ID:          10,
			Source:      0,
			Destination: 1,
			Capacities: infrastructure.Capacity{
				Bands: []infrastructure.Band{
					{Slots: []bool{true, true, true, true}},
				},
			},
		},
	}
	network := infrastructure.Network{Links: links}

	allocated, _ := FirstFit(0, 1, 2, network, routes, 1, "test-id")

	if allocated {
		t.Fatalf("expected allocation to fail due to lack of capacity")
	}
}

func TestFirstFit_FailContiguity(t *testing.T) {
	routes := connections.Routes{
		Paths: []connections.Path{
			{
				Source:      0,
				Destination: 1,
				PathLinks:   [][]int{{0, 1}},
			},
		},
	}

	// Needs 2 contiguous slots, but only has fragmented slots
	links := []infrastructure.Link{
		{
			ID:          10,
			Source:      0,
			Destination: 1,
			Capacities: infrastructure.Capacity{
				Bands: []infrastructure.Band{
					{Slots: []bool{false, true, false, true}},
				},
			},
		},
	}
	network := infrastructure.Network{Links: links}

	allocated, _ := FirstFit(0, 1, 2, network, routes, 1, "test-id")

	if allocated {
		t.Fatalf("expected allocation to fail due to lack of contiguous capacity")
	}
}
