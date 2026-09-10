package allocator

import (
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

// Allocator receives the whole BitRate magnitude selected for the connection
// (every modulation it supports) so it can pick the modulation itself, based
// on each candidate route's length, mirroring the Python reference allocator.
type Allocator func(source, destination int, bitRate connections.BitRate, network infrastructure.Network, path connections.Routes, numberOfBands int, id string, addConnection func(connections.Connection)) bool

func FirstFit(source int, destination int, bitRate connections.BitRate, network infrastructure.Network, path connections.Routes, numberOfBands int, id string, addConnection func(connections.Connection)) bool {

	paths := path.GetPaths(source, destination)

	for _, pathSelected := range paths {
		links := network.GetLinkByPath(pathSelected)
		if len(links) == 0 {
			continue
		}

		// Distance-adaptive modulation selection: pick the modulation with
		// the smallest reach that still covers this route's length.
		length := float64(network.GetPathDistance(links))
		modulation := bitRate.GetDistanceAdaptive(length)
		if modulation == -1 {
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

			bandName := ""
			if band < len(links[0].Capacities.Bands) {
				bandName = links[0].Capacities.Bands[band].Name
			}
			slotCount := bitRate.SlotsForBand(modulation, bandName)
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
							return false
						}
					}

					connection := connections.Connection{
						Id:           id,
						Source:       source,
						Destination:  destination,
						InitialSlot:  currentSlotIndex,
						FinalSlot:    currentSlotIndex + slotCount - 1,
						Slots:        slotCount,
						BandSelected: band,
						Links:        links,
						Allocated:    true,
					}

					if addConnection != nil {
						addConnection(connection)
					}

					return true
				}
			}
		}
	}

	return false
}
