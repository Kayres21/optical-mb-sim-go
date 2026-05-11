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

func (c *Controller) ReleaseConnection(connectionId string) error {
	con, ok := c.Connections[connectionId]
	if !ok {
		return fmt.Errorf("connection with ID %s not found", connectionId)
	}

	for _, link := range con.Links {
		link.ReleaseConnection(con.InitialSlot, con.Slots, con.BandSelected)
	}

	delete(c.Connections, connectionId)
	return nil
}

// ConnectionAllocation delegates to the configured Allocator using the
// controller's own Network and Routes — callers no longer need to pass them in.
func (c *Controller) ConnectionAllocation(source, destination int, getSlot func(band int) int, numberOfBands int, id string) (bool, connections.Connection) {
	return c.Allocator(source, destination, getSlot, c.Network, c.Routes, numberOfBands, id)
}
