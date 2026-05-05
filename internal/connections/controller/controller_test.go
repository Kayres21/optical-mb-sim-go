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

	err := c.ReleaseConnection("1")
	if err != nil {
		t.Errorf("expected no error on release, got %v", err)
	}

	if len(c.Connections) != 0 {
		t.Errorf("expected 0 connections after release, got %d", len(c.Connections))
	}

	err = c.ReleaseConnection("1")
	if err == nil {
		t.Errorf("expected error when releasing non-existent connection")
	}
}

func TestController_ConnectionAllocation(t *testing.T) {
	// Dummy allocator
	dummyAllocator := func(source, destination int, slot int, network infrastructure.Network, path connections.Routes, numberOfBands int, id string) (bool, connections.Connection) {
		return true, connections.Connection{Id: id, Source: source, Destination: destination}
	}

	c := Controller{
		Allocator: dummyAllocator,
	}

	success, conn := c.ConnectionAllocation(0, 1, 1, 1, "test-id")
	if !success {
		t.Errorf("expected successful allocation")
	}
	if conn.Id != "test-id" {
		t.Errorf("expected connection id test-id, got %v", conn.Id)
	}
}
