package allocator

import (
	"log"
	"simulator/internal/connections"
	"simulator/internal/infrastructure"
)

type Allocator func(source, destination int, bitRate connections.BitRate, network infrastructure.Network, path connections.Routes, numberOfBands int) (bool, connections.Connection)

func FirstFit(source int, destination int, bitRate connections.BitRate, network infrastructure.Network, path connections.Routes, numberOfBands int) (bool, connections.Connection) {
	pathSelected := path.GetKshortestPath(0, source, destination)

	links := network.GetLinkByIDs(pathSelected)
	for band := range numberOfBands {
		capacityTotal := make([]bool, len(links[0].GetSlotsByBand(band)))
		log.Printf("Checking band %d with %d links with capacity %d", band, len(links), len(capacityTotal))
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

				return true, connections.Connection{
					Id:         connections.GenerateConnectionID(),
					Source:     source,
					Destination: destination,
					Bitrate:    bitRate.Modulation,
					Links:      links,
					Slots:      len(bitRate.Slots),
					InitialSlot: currentSlotIndex,
					FinalSlot:   currentSlotIndex + len(bitRate.Slots) - 1,
					BandSelected: band,
					Allocated:   true,
			}
		}

	}

	connection := connections.Connection{}

	return true, connection // Return true if allocation is successful
}
