package examples

import (
	"fmt"

	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/defragmentator"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

func main() {
	capacity := infrastructure.Capacity{
		Bands: []infrastructure.Band{{
			Name:     "A",
			SlotsLen: 4,
			Slots:    make([]bool, 4),
		}},
	}

	network := infrastructure.Network{
		Name:  "example",
		Alias: "example",
		Nodes: []infrastructure.Node{{ID: 0}, {ID: 1}},
		Links: []infrastructure.Link{{
			ID:          0,
			Source:      0,
			Destination: 1,
			Length:      1,
			Capacities:  capacity,
		}},
	}

	link := &network.Links[0]
	if err := link.AssignConnection(0, 1, 0); err != nil {
		panic(err)
	}
	if err := link.AssignConnection(2, 1, 0); err != nil {
		panic(err)
	}

	activeConnections := map[string]connections.Connection{
		"1": {
			Id:           "1",
			Source:       0,
			Destination:  1,
			Links:        []*infrastructure.Link{link},
			Slots:        1,
			InitialSlot:  0,
			FinalSlot:    0,
			BandSelected: 0,
			Allocated:    true,
		},
		"2": {
			Id:           "2",
			Source:       0,
			Destination:  1,
			Links:        []*infrastructure.Link{link},
			Slots:        1,
			InitialSlot:  2,
			FinalSlot:    2,
			BandSelected: 0,
			Allocated:    true,
		},
	}

	fmt.Println("before defrag:", link.GetSlotsByBand(0))
	fmt.Println("decision:", defragmentator.DefaultDecision(network, activeConnections, connections.ConnectionEvent{}, 1))

	if _, err := defragmentator.FirstFitActiveConnections(network, activeConnections, connections.Routes{}, 1); err != nil {
		panic(err)
	}

	fmt.Println("after defrag:", link.GetSlotsByBand(0))
	fmt.Printf("connection 1 => slot %d\n", activeConnections["1"].InitialSlot)
	fmt.Printf("connection 2 => slot %d\n", activeConnections["2"].InitialSlot)

	if _, err := defragmentator.DefaultAction(network, activeConnections, connections.Routes{}, nil, 1); err != nil {
		panic(err)
	}
}
