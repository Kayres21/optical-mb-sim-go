package defragmentator

import (
	"fmt"
	"sort"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

type DecisionFunc func(network infrastructure.Network, connections map[string]connections.Connection, event connections.ConnectionEvent, numberOfBands int) bool

// ActionFunc returns the number of connections that were actually relocated
// to a better slot/band, plus any error encountered.
type ActionFunc func(network infrastructure.Network, connections map[string]connections.Connection, routes connections.Routes, alloc allocator.Allocator, numberOfBands int) (int, error)

const (
	DefragNone          = "none"
	DefragBeforeArrival = "before_arrival"
	DefragAfterBlock    = "after_block"
	DefragAfterAssign   = "after_assign"
)

type Defragmenter struct {
	Network       infrastructure.Network
	Connections   map[string]connections.Connection
	Routes        connections.Routes
	Allocator     allocator.Allocator
	NumberOfBands int
	PendingEvents []connections.ConnectionEvent
	nextIndex     int
	finished      bool
}

func NewDefragmenter(network infrastructure.Network, connectionsMap map[string]connections.Connection, routes connections.Routes, alloc allocator.Allocator, numberOfBands int, events []connections.ConnectionEvent) *Defragmenter {
	clonedConnections := make(map[string]connections.Connection, len(connectionsMap))
	for id, conn := range connectionsMap {
		clonedConnections[id] = conn
	}

	clonedEvents := make([]connections.ConnectionEvent, len(events))
	copy(clonedEvents, events)

	return &Defragmenter{
		Network:       network,
		Connections:   clonedConnections,
		Routes:        routes,
		Allocator:     alloc,
		NumberOfBands: numberOfBands,
		PendingEvents: clonedEvents,
		nextIndex:     0,
		finished:      false,
	}
}

func (d *Defragmenter) Process() {
	d.finished = true
}

func DefaultDecision(network infrastructure.Network, connections map[string]connections.Connection, event connections.ConnectionEvent, numberOfBands int) bool {
	return true
}

func DefaultAction(network infrastructure.Network, connectionsMap map[string]connections.Connection, routes connections.Routes, alloc allocator.Allocator, numberOfBands int) (int, error) {
	return 0, nil
}

func BeforeArrivalDecision(network infrastructure.Network, connections map[string]connections.Connection, event connections.ConnectionEvent, numberOfBands int) bool {
	return DefaultDecision(network, connections, event, numberOfBands)
}

func BeforeArrivalAction(network infrastructure.Network, connectionsMap map[string]connections.Connection, routes connections.Routes, alloc allocator.Allocator, numberOfBands int) (int, error) {
	return DefaultAction(network, connectionsMap, routes, alloc, numberOfBands)
}

// FirstFitActiveConnections attempts to compact every active connection into
// a lower-fragmentation slot/band, returning how many connections actually moved.
func FirstFitActiveConnections(network infrastructure.Network, activeConnections map[string]connections.Connection, routes connections.Routes, numberOfBands int) (int, error) {

	if len(activeConnections) == 0 {
		return 0, nil
	}

	moved := 0
	ids := make([]string, 0, len(activeConnections))
	for id := range activeConnections {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		conn := activeConnections[id]
		if !conn.Allocated || conn.Slots <= 0 || len(conn.Links) == 0 {
			continue
		}

		originalSlot := conn.InitialSlot
		originalBand := conn.BandSelected
		routeLinks := conn.Links
		if len(routeLinks) == 0 {
			continue
		}

		// Links keep FragmentationRatioByBand up to date on every assign/release,
		// so the current ratio can be read directly while the connection is still in place.
		currentRatio := averageFragmentationRatio(routeLinks, originalBand)

		for _, link := range routeLinks {
			if err := link.ReleaseConnection(originalSlot, conn.Slots, originalBand); err != nil {
				return moved, fmt.Errorf("release connection %s before reallocation: %w", id, err)
			}
		}

		placed := false
		for band := 0; band < numberOfBands && band < len(routeLinks[0].Capacities.Bands); band++ {
			// Only candidate slots earlier than the connection's original position count as an improvement.
			candidateStarts, ok := findFirstFitRouteSlot(routeLinks, band, conn.Slots, originalSlot)
			if !ok {
				continue
			}

			for _, start := range candidateStarts {
				for _, link := range routeLinks {
					if err := link.AssignConnection(start, conn.Slots, band); err != nil {
						for _, restoreLink := range routeLinks {
							_ = restoreLink.AssignConnection(originalSlot, conn.Slots, originalBand)
						}
						return moved, fmt.Errorf("assign connection %s in route: %w", id, err)
					}
				}

				candidateRatio := averageFragmentationRatio(routeLinks, band)
				if candidateRatio >= currentRatio {
					for _, link := range routeLinks {
						if err := link.ReleaseConnection(start, conn.Slots, band); err != nil {
							return moved, fmt.Errorf("revert candidate placement for connection %s: %w", id, err)
						}
					}
					continue
				}

				conn.InitialSlot = start
				conn.FinalSlot = start + conn.Slots - 1
				conn.BandSelected = band
				conn.Allocated = true
				activeConnections[id] = conn
				placed = true
				moved++
				break
			}
			if placed {
				break
			}
		}

		if !placed {
			for _, link := range routeLinks {
				if err := link.AssignConnection(originalSlot, conn.Slots, originalBand); err != nil {
					return moved, fmt.Errorf("restore connection %s after failed reallocation: %w", id, err)
				}
			}
			conn.InitialSlot = originalSlot
			conn.FinalSlot = originalSlot + conn.Slots - 1
			conn.BandSelected = originalBand
			conn.Allocated = true
			activeConnections[id] = conn
		}
	}

	return moved, nil
}

func averageFragmentationRatio(routeLinks []*infrastructure.Link, band int) float64 {
	if len(routeLinks) == 0 {
		return 0
	}

	total := 0.0
	for _, link := range routeLinks {
		total += link.GetFragmentationRatioByBand(band)
	}
	return total / float64(len(routeLinks))
}

// findFirstFitRouteSlot returns every free start index below maxStart where
// the whole route has slotCount contiguous free slots on the given band.
func findFirstFitRouteSlot(routeLinks []*infrastructure.Link, band, slotCount, maxStart int) ([]int, bool) {
	if len(routeLinks) == 0 || slotCount <= 0 {
		return nil, false
	}

	maxLen := len(routeLinks[0].GetSlotsByBand(band))
	lastStart := maxLen - slotCount
	if maxStart-1 < lastStart {
		lastStart = maxStart - 1
	}
	starts := make([]int, 0)
	for start := 0; start <= lastStart; start++ {
		fits := true
		for _, link := range routeLinks {
			if band >= len(link.Capacities.Bands) {
				fits = false
				break
			}
			bandSlots := link.GetSlotsByBand(band)
			if len(bandSlots) < start+slotCount {
				fits = false
				break
			}
			for idx := start; idx < start+slotCount; idx++ {
				if bandSlots[idx] {
					fits = false
					break
				}
			}
			if !fits {
				break
			}
		}
		if fits {
			starts = append(starts, start)
		}
	}
	if len(starts) == 0 {
		return nil, false
	}
	return starts, true
}
