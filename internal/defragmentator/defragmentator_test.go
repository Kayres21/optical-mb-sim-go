package defragmentator

import (
	"testing"

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

func TestDefaultDecisionAlwaysTrue(t *testing.T) {
	network := makeTestNetwork()
	if !DefaultDecision(network, nil, connections.ConnectionEvent{}, 1) {
		t.Fatal("expected default decision to always return true")
	}
}

func TestDefaultActionDoesNothing(t *testing.T) {
	network := makeTestNetwork()
	connectionsMap := map[string]connections.Connection{}
	if _, err := DefaultAction(network, connectionsMap, connections.Routes{}, nil, 1); err != nil {
		t.Fatalf("expected default action to be a no-op, got: %v", err)
	}
}

func TestFirstFitActiveConnectionsCompactsRoute(t *testing.T) {
	network := makeTestNetwork()
	link := &network.Links[0]
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

	if moved, err := FirstFitActiveConnections(network, connectionsMap, connections.Routes{}, 1); err != nil {
		t.Fatalf("first-fit compaction failed: %v", err)
	} else if moved != 1 {
		t.Fatalf("expected 1 connection to be moved, got %d", moved)
	}

	slots := link.GetSlotsByBand(0)
	if got := slots[0]; !got {
		t.Fatalf("expected slot 0 to remain occupied, got %v", slots)
	}
	if got := slots[1]; !got {
		t.Fatalf("expected slot 1 to be occupied after moving the later connection forward, got %v", slots)
	}
	if got := slots[2]; got {
		t.Fatalf("expected slot 2 to be free after compaction, got %v", slots)
	}
	if got := slots[3]; got {
		t.Fatalf("expected slot 3 to stay free, got %v", slots)
	}

	if conn1 := connectionsMap["1"]; conn1.InitialSlot != 0 {
		t.Fatalf("expected connection 1 to stay at slot 0, got %+v", conn1)
	}
	if conn2 := connectionsMap["2"]; conn2.InitialSlot != 1 {
		t.Fatalf("expected connection 2 to move to slot 1 because only smaller indexes are considered, got %+v", conn2)
	}
}
