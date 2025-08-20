package allocator

import (
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

type Allocator func(source, destination int, bitRate connections.BitRate, network infrastructure.Network, path connections.Routes, numberOfBands int) (bool, connections.Connection)

func FirstFit(source int, destination int, bitRate connections.BitRate, network infrastructure.Network, path connections.Routes, numberOfBands int) (bool, connections.Connection) {
	pathSelected := path.GetKshortestPath(0, source, destination)

	links := network.GetLinkByIDs(pathSelected)
	for band := range numberOfBands {
		capacityTotal := make([]bool, len(links[0].GetSlotsByBand(band)))
		for _, link := range links {

			capacity := link.GetSlotsByBand(band)

			if len(capacity) != len(capacityTotal) {
				return false, connections.Connection{}
			}

			for i := range capacity {
				capacityTotal[i] = capacityTotal[i] || capacity[i]
			}

		}

		continousSlots := 0
		currentSlotIndex := 0
		for i := range capacityTotal {
			if !capacityTotal[i] {
				continousSlots++
			} else {
				continousSlots = 0
				currentSlotIndex = i + 1
			}

			if continousSlots == len(bitRate.Slots) {

				for _, link := range links {

					link.AssignConnection(currentSlotIndex, len(bitRate.Slots), band)

				}
				linksVal := make([]infrastructure.Link, len(links))
				for i, l := range links {
					linksVal[i] = *l
				}
				return true, connections.Connection{
					Id:           connections.GenerateConnectionID(),
					Source:       source,
					Destination:  destination,
					BitRate:      bitRate.Modulation,
					Links:        linksVal,
					Slots:        len(bitRate.Slots),
					InitialSlot:  currentSlotIndex,
					FinalSlot:    currentSlotIndex + len(bitRate.Slots) - 1,
					BandSelected: band,
					Allocated:    true,
				}
			}
		}

	}

	connection := connections.Connection{}

	return false, connection
}
