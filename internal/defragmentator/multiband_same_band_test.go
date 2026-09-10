package defragmentator

import (
	"testing"

	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

func TestMultiBandSameBandActiveConnections_DoesNotMigrateBands(t *testing.T) {
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

	moved, err := MultiBandSameBandActiveConnections(network, activeConnections, connections.Routes{}, 2)
	if err != nil {
		t.Fatalf("same-band compaction failed: %v", err)
	}
	if moved != 0 {
		t.Fatalf("expected no move because only another band improves FR, got %d moves", moved)
	}

	connection := activeConnections["target"]
	if connection.BandSelected != 0 || connection.InitialSlot != 1 {
		t.Fatalf("expected connection to remain in band 0 at slot 1, got band %d slot %d", connection.BandSelected, connection.InitialSlot)
	}
	if !link.GetSlotsByBand(0)[1] {
		t.Fatal("expected original slot to remain occupied in band 0")
	}
	if link.GetSlotsByBand(1)[0] {
		t.Fatal("expected band 1 to remain unused")
	}
}
