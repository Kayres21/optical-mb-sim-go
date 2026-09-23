package simulations

import (
	"testing"

	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

func TestSweepSeedsAreIndependentAcrossIterations(t *testing.T) {
	first := sweepSeedsForIteration(0)
	second := sweepSeedsForIteration(1)

	if first == second {
		t.Fatal("expected different seed sets for different sweep iterations")
	}
	if first.arrive == second.arrive || first.departure == second.departure || first.bitrate == second.bitrate || first.source == second.source || first.destination == second.destination || first.band == second.band {
		t.Fatal("expected each sweep seed component to change across iterations")
	}
}

func TestSweepRunnerClonesNetworkPerIteration(t *testing.T) {
	network := infrastructure.Network{
		Name:  "Test",
		Alias: "T",
		Links: []infrastructure.Link{{
			ID:          1,
			Source:      1,
			Destination: 2,
			Capacities: infrastructure.Capacity{Bands: []infrastructure.Band{{
				Name:     "C",
				SlotsLen: 4,
				Slots:    []bool{false, false, false, false},
			}}},
			FragmentationRatioByBand: []float64{0},
		}},
	}

	clone := network.Clone()
	clone.Links[0].Capacities.Bands[0].Slots[0] = true
	if network.Links[0].Capacities.Bands[0].Slots[0] {
		t.Fatal("network clone mutated original state")
	}
}
