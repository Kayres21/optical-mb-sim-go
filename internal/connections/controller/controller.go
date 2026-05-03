package controller

import (
	"fmt"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

type Controller struct {
	Routes      connections.Routes
	Connections []connections.Connection
	Network     infrastructure.Network
	Allocator   allocator.Allocator
}



func (c *Controller) AddConnection(connection connections.Connection) {
	c.Connections = append(c.Connections, connection)
}

func New(pathToRoutes string, network infrastructure.Network, allocator allocator.Allocator) (Controller, error) {
	routes, err := connections.ReadRoutesFile(pathToRoutes)
	if err != nil {
		return Controller{}, err
	}

	return Controller{
		Routes:      routes,
		Connections: []connections.Connection{},
		Network:     network,
		Allocator:   allocator,
	}, nil
}

func (c *Controller) GetConnectionById(id string) (connections.Connection, bool) {
	for _, connection := range c.Connections {
		if connection.Id == id {
			return connection, true
		}
	}
	return connections.Connection{}, false
}

func (c *Controller) ReleaseConnection(connectionId string) error {

	con, valid := c.GetConnectionById(connectionId)

	if !valid {
		return fmt.Errorf("connection with ID %s not found", connectionId)
	}

	links := con.Links

	for _, link := range links {
		link.ReleaseConnection(con.InitialSlot, con.Slots, con.BandSelected)

	}

	for i, connection := range c.Connections {
		if connection.Id == connectionId {
			c.Connections = append(c.Connections[:i], c.Connections[i+1:]...)
			break
		}
	}

	return nil
}

func (c *Controller) ConectionAllocation(source, destination int, slot int, network infrastructure.Network, path connections.Routes, numberOfBands int, id string) (bool, connections.Connection) {
	return c.Allocator(source, destination, slot, network, path, numberOfBands, id)
}
