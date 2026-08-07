package allocator

import (
	"os"
	"path/filepath"
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

	allocated := FirstFit(0, 1, func(int) int { return 2 }, network, routes, 1, "test-id", nil)

	if !allocated {
		t.Fatalf("expected allocation to succeed")
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

	allocated := FirstFit(0, 1, func(int) int { return 2 }, network, routes, 1, "test-id", nil)

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

	allocated := FirstFit(0, 1, func(int) int { return 2 }, network, routes, 1, "test-id", nil)

	if allocated {
		t.Fatalf("expected allocation to fail due to lack of contiguous capacity")
	}
}

func TestFirstFitFromFile_DoesNotAllocateConnections(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "allocations.csv")
	csvContents := "id,source,destination,decision\nreq-1,0,1,ALLOCATED\n"
	if err := os.WriteFile(csvPath, []byte(csvContents), 0o644); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	alloc, err := FirstFitFromFile(csvPath)
	if err != nil {
		t.Fatalf("create allocator: %v", err)
	}

	routes := connections.Routes{
		Paths: []connections.Path{{
			Source:      0,
			Destination: 1,
			PathLinks:   [][]int{{0, 1}},
		}},
	}

	links := []infrastructure.Link{{
		ID:          10,
		Source:      0,
		Destination: 1,
		Capacities:  infrastructure.Capacity{Bands: []infrastructure.Band{{Slots: []bool{false, false, false, false}}}},
	}}
	network := infrastructure.Network{Links: links}

	allocated := alloc(0, 1, func(int) int { return 2 }, network, routes, 1, "req-1", nil)
	if !allocated {
		t.Fatalf("expected allocation to be allowed")
	}

	if links[0].Capacities.Bands[0].Slots[0] {
		t.Fatalf("expected no slot to be marked allocated")
	}
}

func TestFirstFitFromFile_UsesCSVDecisionEvenWhenCapacityLooksFull(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "allocations.csv")
	csvContents := "id,source,destination,decision\nreq-1,0,1,ALLOCATED\n"
	if err := os.WriteFile(csvPath, []byte(csvContents), 0o644); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	alloc, err := FirstFitFromFile(csvPath)
	if err != nil {
		t.Fatalf("create allocator: %v", err)
	}

	routes := connections.Routes{
		Paths: []connections.Path{{
			Source:      0,
			Destination: 1,
			PathLinks:   [][]int{{0, 1}},
		}},
	}

	links := []infrastructure.Link{{
		ID:          10,
		Source:      0,
		Destination: 1,
		Capacities:  infrastructure.Capacity{Bands: []infrastructure.Band{{Slots: []bool{true, true, true, true}}}},
	}}
	network := infrastructure.Network{Links: links}

	allocated := alloc(0, 1, func(int) int { return 1 }, network, routes, 1, "req-1", nil)
	if !allocated {
		t.Fatalf("expected ALLOCATED CSV entry to be accepted")
	}
}

func TestFirstFit_SearchesAllRoutes(t *testing.T) {
	routes := connections.Routes{
		Paths: []connections.Path{
			{
				Source:      0,
				Destination: 1,
				PathLinks:   [][]int{{0, 2, 1}, {0, 3, 1}},
			},
		},
	}

	links := []infrastructure.Link{
		{
			ID:          10,
			Source:      0,
			Destination: 2,
			Capacities:  infrastructure.Capacity{Bands: []infrastructure.Band{{Slots: []bool{false, true, false, true}}}},
		},
		{
			ID:          20,
			Source:      2,
			Destination: 1,
			Capacities:  infrastructure.Capacity{Bands: []infrastructure.Band{{Slots: []bool{false, true, false, true}}}},
		},
		{
			ID:          30,
			Source:      0,
			Destination: 3,
			Capacities:  infrastructure.Capacity{Bands: []infrastructure.Band{{Slots: []bool{false, false, false, false}}}},
		},
		{
			ID:          40,
			Source:      3,
			Destination: 1,
			Capacities:  infrastructure.Capacity{Bands: []infrastructure.Band{{Slots: []bool{false, false, false, false}}}},
		},
	}
	network := infrastructure.Network{Links: links}

	allocated := FirstFit(0, 1, func(int) int { return 2 }, network, routes, 1, "test-id", nil)

	if !allocated {
		t.Fatalf("expected allocation to succeed on the second route")
	}
}
