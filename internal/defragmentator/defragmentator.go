package defragmentator

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

type DecisionFunc func(network infrastructure.Network, connections map[string]connections.Connection, event connections.ConnectionEvent, numberOfBands int) bool

type ActionFunc func(network infrastructure.Network, connections map[string]connections.Connection, routes connections.Routes, alloc allocator.Allocator, numberOfBands int) error

const (
	DefragNone          = "none"
	DefragBeforeArrival = "before_arrival"
	DefragAfterBlock    = "after_block"
	DefragAfterAssign   = "after_assign"
)

func DefaultDecision(network infrastructure.Network, connections map[string]connections.Connection, event connections.ConnectionEvent, numberOfBands int) bool {
	for _, link := range network.Links {
		for band := 0; band < numberOfBands && band < len(link.Capacities.Bands); band++ {
			slots := link.GetSlotsByBand(band)
			used := false
			for _, occupied := range slots {
				if occupied {
					used = true
				} else if used {
					return true
				}
			}
		}
	}
	return false
}

func DefaultAction(network infrastructure.Network, connectionsMap map[string]connections.Connection, routes connections.Routes, alloc allocator.Allocator, numberOfBands int) error {
	if len(connectionsMap) == 0 {
		return nil
	}

	type connEntry struct {
		id   string
		conn connections.Connection
		idx  int
	}

	entries := make([]connEntry, 0, len(connectionsMap))
	for id, conn := range connectionsMap {
		idx, err := strconv.Atoi(id)
		if err != nil {
			idx = int(^uint(0) >> 1)
		}
		entries = append(entries, connEntry{id: id, conn: conn, idx: idx})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].idx != entries[j].idx {
			return entries[i].idx < entries[j].idx
		}
		return entries[i].id < entries[j].id
	})

	original := make(map[string]connections.Connection, len(connectionsMap))
	for id, conn := range connectionsMap {
		original[id] = conn
	}

	for _, conn := range original {
		for _, link := range conn.Links {
			_ = link.ReleaseConnection(conn.InitialSlot, conn.Slots, conn.BandSelected)
		}
	}

	updated := make(map[string]connections.Connection, len(original))
	for _, entry := range entries {
		connection := entry.conn
		getSlot := func(band int) int {
			if band == connection.BandSelected {
				return connection.Slots
			}
			return 0
		}

		assigned, newConn := alloc(connection.Source, connection.Destination, getSlot, network, routes, numberOfBands, connection.Id)
		if !assigned {
			for _, conn := range original {
				for _, link := range conn.Links {
					_ = link.AssignConnection(conn.InitialSlot, conn.Slots, conn.BandSelected)
				}
			}
			for id, conn := range original {
				connectionsMap[id] = conn
			}
			return fmt.Errorf("failed to reallocate connection %s during defragmentation", entry.id)
		}

		updated[entry.id] = newConn
	}

	for id, conn := range updated {
		connectionsMap[id] = conn
	}

	return nil
}
