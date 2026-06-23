package defragmentator

import (
	"testing"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

func makeTestNetwork() infrastructure.Network {
	capacity := infrastructure.Capacity{
		Bands: []infrastructure.Band{{
			Name:     "A",
			SlotsLen: 4,
			Slots:    make([]bool, 4),
		}},
	}

	return infrastructure.Network{
		Name:  "test",
		Alias: "test",
		Nodes: []infrastructure.Node{{ID: 0}, {ID: 1}},
		Links: []infrastructure.Link{{
			ID:          0,
			Source:      0,
			Destination: 1,
			Length:      1,
			Capacities:  capacity,
		}},
	}
}

func makeTestRoutes() connections.Routes {
	return connections.Routes{
		Paths: []connections.Path{{
			Source:      0,
			Destination: 1,
			PathLinks:   [][]int{{0, 1}},
		}},
	}
}

func TestDefaultDecisionDetectsFragmentation(t *testing.T) {
	network := makeTestNetwork()
	link := network.GetLinkByPath([]int{0, 1})[0]

	if err := link.AssignConnection(0, 1, 0); err != nil {
		t.Fatal(err)
	}
	if err := link.AssignConnection(2, 1, 0); err != nil {
		t.Fatal(err)
	}

	connectionsMap := map[string]connections.Connection{
		"1": {Id: "1", Source: 0, Destination: 1, Links: []*infrastructure.Link{link}, Slots: 1, InitialSlot: 0, FinalSlot: 0, BandSelected: 0, Allocated: true},
		"2": {Id: "2", Source: 0, Destination: 1, Links: []*infrastructure.Link{link}, Slots: 1, InitialSlot: 2, FinalSlot: 2, BandSelected: 0, Allocated: true},
	}

	event := connections.ConnectionEvent{Id: "1", Source: 0, Destination: 1}
	if !DefaultDecision(network, connectionsMap, event, 1) {
		t.Fatalf("expected default decision to detect fragmentation")
	}
}

func TestDefaultActionReallocatesConnections(t *testing.T) {
	network := makeTestNetwork()
	routes := makeTestRoutes()
	link := network.GetLinkByPath([]int{0, 1})[0]

	if err := link.AssignConnection(0, 1, 0); err != nil {
		t.Fatal(err)
	}
	if err := link.AssignConnection(2, 1, 0); err != nil {
		t.Fatal(err)
	}

	connectionsMap := map[string]connections.Connection{
		"1": {Id: "1", Source: 0, Destination: 1, Links: []*infrastructure.Link{link}, Slots: 1, InitialSlot: 0, FinalSlot: 0, BandSelected: 0, Allocated: true},
		"2": {Id: "2", Source: 0, Destination: 1, Links: []*infrastructure.Link{link}, Slots: 1, InitialSlot: 2, FinalSlot: 2, BandSelected: 0, Allocated: true},
	}

	if err := DefaultAction(network, connectionsMap, routes, allocator.FirstFit, 1); err != nil {
		t.Fatalf("defragmentation action failed: %v", err)
	}

	slots := link.GetSlotsByBand(0)
	expected := []bool{true, true, false, false}
	for i, got := range slots {
		if got != expected[i] {
			t.Fatalf("expected slots %v, got %v", expected, slots)
		}
	}

	conn1, ok := connectionsMap["1"]
	if !ok || conn1.InitialSlot != 0 {
		t.Fatalf("expected connection 1 to be reallocated to slot 0, got %+v", conn1)
	}

	conn2, ok := connectionsMap["2"]
	if !ok || conn2.InitialSlot != 1 {
		t.Fatalf("expected connection 2 to be reallocated to slot 1, got %+v", conn2)
	}
}
