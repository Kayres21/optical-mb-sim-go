package defragmentator

import (
	"fmt"

	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

// MultiBandSameBandActiveConnections compacts active connections only within
// their original band; it never migrates a connection to another band.
func MultiBandSameBandActiveConnections(network infrastructure.Network, activeConnections map[string]connections.Connection, routes connections.Routes, numberOfBands int) (int, error) {
	if len(activeConnections) == 0 || numberOfBands <= 0 {
		return 0, nil
	}

	moved := 0
	for _, id := range sortedConnectionIDs(activeConnections) {
		connection := activeConnections[id]
		if !connection.Allocated || connection.Slots <= 0 || len(connection.Links) == 0 {
			continue
		}

		originalSlot := connection.InitialSlot
		originalBand := connection.BandSelected
		if originalBand < 0 || originalBand >= numberOfBands {
			continue
		}

		routeLinks := connection.Links
		currentRatio := averageFragmentationRatio(routeLinks, originalBand)
		if err := releaseRoute(routeLinks, originalSlot, connection.Slots, originalBand); err != nil {
			return moved, fmt.Errorf("release connection %s before reallocation: %w", id, err)
		}

		placed := false
		candidateStarts, ok := findFirstFitRouteSlot(routeLinks, originalBand, connection.Slots, originalSlot)
		if ok {
			for _, start := range candidateStarts {
				if err := assignRoute(routeLinks, start, connection.Slots, originalBand); err != nil {
					restoreRoute(routeLinks, originalSlot, connection.Slots, originalBand)
					return moved, fmt.Errorf("assign connection %s in original band: %w", id, err)
				}

				if averageFragmentationRatio(routeLinks, originalBand) >= currentRatio {
					if err := releaseRoute(routeLinks, start, connection.Slots, originalBand); err != nil {
						return moved, fmt.Errorf("revert candidate placement for connection %s: %w", id, err)
					}
					continue
				}

				connection.InitialSlot = start
				connection.FinalSlot = start + connection.Slots - 1
				activeConnections[id] = connection
				moved++
				placed = true
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
