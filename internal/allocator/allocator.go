package allocator

import (
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

func countTrue(s []bool) int {
	count := 0
	for _, v := range s {
		if v {
			count++
		}
	}
	return count
}

type Allocator func(source, destination int, slot int, network infrastructure.Network, path connections.Routes, numberOfBands int, id string) (bool, connections.Connection)

func FirstFit(source int, destination int, slot int, network infrastructure.Network, path connections.Routes, numberOfBands int, id string) (bool, connections.Connection) {

	pathSelected := path.GetKshortestPath(0, source, destination)

	links := network.GetLinkByPath(pathSelected)

	for band := range numberOfBands {
		capacityTotal := make([]bool, len(links[0].GetSlotsByBand(band)))

		for _, link := range links {

			capacity := link.GetSlotsByBand(band)

			//fmt.Println(countTrue(capacity), "occupied slots in link", link.ID, "band", band)

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

			if continousSlots == slot {

				for _, link := range links {

					link.AssignConnection(currentSlotIndex, slot, band)

				}
				return true, connections.Connection{
					Id:           id,
					Source:       source,
					Destination:  destination,
					Links:        links,
					Slots:        slot,
					InitialSlot:  currentSlotIndex,
					FinalSlot:    currentSlotIndex + slot - 1,
					BandSelected: band,
					Allocated:    true,
				}
			}
		}

	}

	connection := connections.Connection{}

	return false, connection
}
