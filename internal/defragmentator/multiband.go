package defragmentator

import (
	"fmt"
	"sort"

	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

// MultiBandActiveConnections first tries to compact each connection in its
// current band, then tries the remaining enabled bands.
func MultiBandActiveConnections(network infrastructure.Network, activeConnections map[string]connections.Connection, routes connections.Routes, numberOfBands int) (int, error) {
	if len(activeConnections) == 0 || numberOfBands <= 0 {
		return 0, nil
	}

	moved := 0
	ids := sortedConnectionIDs(activeConnections)
	for _, id := range ids {
		connection := activeConnections[id]
		if !connection.Allocated || connection.Slots <= 0 || len(connection.Links) == 0 {
			continue
		}

		originalSlot := connection.InitialSlot
		originalBand := connection.BandSelected
		routeLinks := connection.Links
		if originalBand < 0 || originalBand >= numberOfBands {
			continue
		}

		currentRatio := averageFragmentationRatio(routeLinks, originalBand)
		for _, link := range routeLinks {
			if err := link.ReleaseConnection(originalSlot, connection.Slots, originalBand); err != nil {
				return moved, fmt.Errorf("release connection %s before reallocation: %w", id, err)
			}
		}

		placed := false
		for _, band := range orderedBands(originalBand, numberOfBands) {
			candidateStarts, ok := findFirstFitRouteSlot(routeLinks, band, connection.Slots, originalSlot)
			if !ok {
				continue
			}

			for _, start := range candidateStarts {
				if err := assignRoute(routeLinks, start, connection.Slots, band); err != nil {
					restoreRoute(routeLinks, originalSlot, connection.Slots, originalBand)
					return moved, fmt.Errorf("assign connection %s in route: %w", id, err)
				}

				if averageFragmentationRatio(routeLinks, band) >= currentRatio {
					if err := releaseRoute(routeLinks, start, connection.Slots, band); err != nil {
						return moved, fmt.Errorf("revert candidate placement for connection %s: %w", id, err)
					}
					continue
				}

				connection.InitialSlot = start
				connection.FinalSlot = start + connection.Slots - 1
				connection.BandSelected = band
				activeConnections[id] = connection
				moved++
				placed = true
				break
			}
			if placed {
				break
			}
		}

		if !placed {
			if err := assignRoute(routeLinks, originalSlot, connection.Slots, originalBand); err != nil {
				return moved, fmt.Errorf("restore connection %s after failed reallocation: %w", id, err)
			}
		}
	}

	return moved, nil
}

func sortedConnectionIDs(activeConnections map[string]connections.Connection) []string {
	ids := make([]string, 0, len(activeConnections))
	for id := range activeConnections {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func orderedBands(originalBand, numberOfBands int) []int {
	bands := make([]int, 0, numberOfBands)
	bands = append(bands, originalBand)
	for band := 0; band < numberOfBands; band++ {
		if band != originalBand {
			bands = append(bands, band)
		}
	}
	return bands
}

func assignRoute(links []*infrastructure.Link, initialSlot, slotCount, band int) error {
	for _, link := range links {
		if err := link.AssignConnection(initialSlot, slotCount, band); err != nil {
			return err
		}
	}
	return nil
}

func releaseRoute(links []*infrastructure.Link, initialSlot, slotCount, band int) error {
	for _, link := range links {
		if err := link.ReleaseConnection(initialSlot, slotCount, band); err != nil {
			return err
		}
	}
	return nil
}

func restoreRoute(links []*infrastructure.Link, initialSlot, slotCount, band int) {
	for _, link := range links {
		_ = link.AssignConnection(initialSlot, slotCount, band)
	}
}
