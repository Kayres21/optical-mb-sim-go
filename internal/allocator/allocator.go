package allocator

import (
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

type Allocator func(source, destination int, getSlot func(band int) int, network infrastructure.Network, path connections.Routes, numberOfBands int, id string) (bool, connections.Connection)

func FirstFit(source int, destination int, getSlot func(band int) int, network infrastructure.Network, path connections.Routes, numberOfBands int, id string) (bool, connections.Connection) {

	paths := path.GetPaths(source, destination)

	for _, pathSelected := range paths {
		links := network.GetLinkByPath(pathSelected)
		if len(links) == 0 {
			continue
		}

		for band := 0; band < numberOfBands; band++ {
			bandCapacity := links[0].GetSlotsByBand(band)
			capacityTotal := make([]bool, len(bandCapacity))
			validBand := true

			for _, link := range links {
				capacity := link.GetSlotsByBand(band)
				if len(capacity) != len(capacityTotal) {
					validBand = false
					break
				}

				for i := range capacity {
					capacityTotal[i] = capacityTotal[i] || capacity[i]
				}
			}

			if !validBand {
				continue
			}

			slotCount := getSlot(band)
			if slotCount == 0 {
				continue
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

				if continousSlots == slotCount {
					for _, link := range links {
						if err := link.AssignConnection(currentSlotIndex, slotCount, band); err != nil {
							return false, connections.Connection{}
						}
					}

					return true, connections.Connection{
						Id:           id,
						Source:       source,
						Destination:  destination,
						Links:        links,
						Slots:        slotCount,
						InitialSlot:  currentSlotIndex,
						FinalSlot:    currentSlotIndex + slotCount - 1,
						BandSelected: band,
						Allocated:    true,
					}
				}
			}
		}
	}

	return false, connections.Connection{}
}
