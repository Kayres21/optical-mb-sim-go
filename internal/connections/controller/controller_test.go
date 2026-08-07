package controller

import (
	"testing"

	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

func TestController_AddAndGetConnection(t *testing.T) {
	c := Controller{
		Connections: make(map[string]connections.Connection),
	}

	conn := connections.Connection{Id: "1", Source: 0, Destination: 1}
	c.AddConnection(conn)

	if len(c.Connections) != 1 {
		t.Fatalf("expected 1 connection, got %d", len(c.Connections))
	}

	retrieved, found := c.GetConnectionById("1")
	if !found {
		t.Fatalf("expected to find connection")
	}
	if retrieved.Id != "1" {
		t.Errorf("expected connection id 1, got %s", retrieved.Id)
	}

	_, found = c.GetConnectionById("2")
	if found {
		t.Errorf("did not expect to find connection 2")
	}
}

func TestController_ReleaseConnection(t *testing.T) {
	c := Controller{
		Connections: make(map[string]connections.Connection),
	}

	conn := connections.Connection{Id: "1", Source: 0, Destination: 1, Allocated: true}
	c.AddConnection(conn)

	err := c.ReleaseConnection(conn, 0)
	if err != nil {
		t.Errorf("expected no error on release, got %v", err)
	}

	if len(c.Connections) != 0 {
		t.Errorf("expected 0 connections after release, got %d", len(c.Connections))
	}

	err = c.ReleaseConnection(conn, 0)
	if err == nil {
		t.Errorf("expected error when releasing non-existent connection")
	}
}

func TestController_ReleaseConnectionUsesStoredConnectionByID(t *testing.T) {
	c := Controller{
		Connections: make(map[string]connections.Connection),
	}

	link := &infrastructure.Link{
		Capacities: infrastructure.Capacity{Bands: []infrastructure.Band{{Slots: []bool{false, false, false, false, false}}}},
	}
	stored := connections.Connection{
		Id:           "1",
		Source:       0,
		Destination:  1,
		InitialSlot:  1,
		Slots:        2,
		BandSelected: 0,
		Links:        []*infrastructure.Link{link},
	}
	c.AddConnection(stored)

	request := connections.Connection{Id: "1", InitialSlot: 100, Slots: 2, BandSelected: 0}
	if err := c.ReleaseConnection(request, 0); err != nil {
		t.Fatalf("expected release to succeed using stored connection details, got %v", err)
	}
}

func TestController_ConnectionAllocation(t *testing.T) {
	// Dummy allocator
	dummyAllocator := func(source, destination int, getSlot func(int) int, network infrastructure.Network, path connections.Routes, numberOfBands int, id string, addConnection func(connections.Connection)) bool {
		if addConnection != nil {
			addConnection(connections.Connection{Id: id, Source: source, Destination: destination})
		}
		return true
	}

	c := Controller{
		Connections: make(map[string]connections.Connection),
		Allocator:   dummyAllocator,
	}

	success := c.ConnectionAllocation(0, 1, func(int) int { return 1 }, 1, "test-id")
	if !success {
		t.Errorf("expected successful allocation")
	}

	if len(c.Connections) != 1 {
		t.Fatalf("expected connection to be added to controller after allocation")
	}

	stored, found := c.GetConnectionById("test-id")
	if !found {
		t.Fatalf("expected connection to exist in controller")
	}
	if stored.Source != 0 || stored.Destination != 1 {
		t.Errorf("unexpected stored connection details: %+v", stored)
	}
}
