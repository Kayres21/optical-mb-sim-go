package defragmentator

import (
	"testing"

	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

func TestMultiBandActiveConnections_UsesAnotherBandAfterSameBandFails(t *testing.T) {
	network := infrastructure.Network{
		Links: []infrastructure.Link{{
			Capacities: infrastructure.Capacity{Bands: []infrastructure.Band{
				{Name: "C", Slots: make([]bool, 4)},
				{Name: "L", Slots: make([]bool, 4)},
			}},
		}},
	}
	link := &network.Links[0]
	if err := link.AssignConnection(1, 1, 0); err != nil {
		t.Fatal(err)
	}
	if err := link.AssignConnection(3, 1, 0); err != nil {
		t.Fatal(err)
	}

	activeConnections := map[string]connections.Connection{
		"target": {
			Id:           "target",
			Links:        []*infrastructure.Link{link},
			Slots:        1,
			InitialSlot:  1,
			FinalSlot:    1,
			BandSelected: 0,
			Allocated:    true,
		},
	}

	moved, err := MultiBandActiveConnections(network, activeConnections, connections.Routes{}, 2)
	if err != nil {
		t.Fatalf("multiband compaction failed: %v", err)
	}
	if moved != 1 {
		t.Fatalf("expected 1 connection to be moved, got %d", moved)
	}

	connection := activeConnections["target"]
	if connection.BandSelected != 1 || connection.InitialSlot != 0 {
		t.Fatalf("expected fallback to band 1 at slot 0, got band %d slot %d", connection.BandSelected, connection.InitialSlot)
	}
	if link.GetSlotsByBand(0)[1] {
		t.Fatal("expected original slot to be released from band 0")
	}
	if !link.GetSlotsByBand(1)[0] {
		t.Fatal("expected connection to occupy slot 0 in band 1")
	}
}
