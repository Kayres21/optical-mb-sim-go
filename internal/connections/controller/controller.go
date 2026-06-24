package controller

import (
	"fmt"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

type Controller struct {
	Routes      connections.Routes
	Connections map[string]connections.Connection
	Network     infrastructure.Network
	Allocator   allocator.Allocator
	// Optional unassign callback invoked when a connection is released.
	UnassignCallback func(connections.Connection, float64, infrastructure.Network)
	// If true, use MB-style unassign behavior (no callback).
	UseUnassignMB bool
}

func New(pathToRoutes string, network infrastructure.Network, allocator allocator.Allocator) (Controller, error) {
	routes, err := connections.ReadRoutesFile(pathToRoutes)
	if err != nil {
		return Controller{}, err
	}

	return Controller{
		Routes:      routes,
		Connections: make(map[string]connections.Connection),
		Network:     network,
		Allocator:   allocator,
	}, nil
}

func (c *Controller) AddConnection(connection connections.Connection) {
	c.Connections[connection.Id] = connection
}

func (c *Controller) GetConnectionById(id string) (connections.Connection, bool) {
	con, ok := c.Connections[id]
	return con, ok
}

func (c *Controller) GetConnectionByAllocation(source, destination, initialSlot, slotCount, band int) (connections.Connection, bool) {
	for _, con := range c.Connections {
		if con.Source == source && con.Destination == destination && con.InitialSlot == initialSlot && con.Slots == slotCount && con.BandSelected == band {
			return con, true
		}
	}
	return connections.Connection{}, false
}

func (c *Controller) ReleaseConnection(connection connections.Connection, t float64) error {
	if connection.Id != "" {
		if _, ok := c.Connections[connection.Id]; !ok {
			return fmt.Errorf("connection with ID %s not found", connection.Id)
		}
	}

	for _, link := range connection.Links {
		if err := link.ReleaseConnection(connection.InitialSlot, connection.Slots, connection.BandSelected); err != nil {
			return err
		}
	}

	// If there's a callback and MB mode is not enabled, call it
	if c.UnassignCallback != nil && !c.UseUnassignMB {
		c.UnassignCallback(connection, t, c.Network)
	}

	if connection.Id != "" {
		delete(c.Connections, connection.Id)
	}
	return nil
}

// SetUnassignCallback registers a callback invoked when a connection is released.
func (c *Controller) SetUnassignCallback(cb func(connections.Connection, float64, infrastructure.Network)) {
	c.UnassignCallback = cb
	c.UseUnassignMB = false
}

// SetUnassignMB configures controller to use MB unassign behavior (no callback).
func (c *Controller) SetUnassignMB() {
	c.UnassignCallback = nil
	c.UseUnassignMB = true
}

// ConnectionAllocation delegates to the configured Allocator using the
// controller's own Network and Routes — callers no longer need to pass them in.
func (c *Controller) ConnectionAllocation(source, destination int, getSlot func(band int) int, numberOfBands int, id string) (bool, connections.Connection) {
	return c.Allocator(source, destination, getSlot, c.Network, c.Routes, numberOfBands, id)
}
