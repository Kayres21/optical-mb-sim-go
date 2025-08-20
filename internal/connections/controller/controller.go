package controller

import (
	"log"
	"simulator/internal/allocator"
	"simulator/internal/connections"
	"simulator/internal/infrastructure"
)

type Controller struct {
	Routes      connections.Routes
	Connections []connections.Connection
	Network     infrastructure.Network
	Allocator   allocator.Allocator
}

func (c *Controller) GetRoutes() connections.Routes {
	return c.Routes
}

func (c *Controller) GetConnections() []connections.Connection {
	return c.Connections
}

func (c *Controller) GetNetwork() infrastructure.Network {

	return c.Network
}

func (c *Controller) SetRoutes(routes connections.Routes) {
	c.Routes = routes
}

func (c *Controller) SetConnections(connections []connections.Connection) {
	c.Connections = connections
}

func (c *Controller) AddConnection(connection connections.Connection) {
	c.Connections = append(c.Connections, connection)
}

func (c *Controller) SetNetwork(network infrastructure.Network) {
	c.Network = network
}

func (c *Controller) GetConnectionById(id string) (connections.Connection, bool) {
	for _, connection := range c.Connections {
		if connection.Id == id {
			return connection, true
		}
	}
	return connections.Connection{}, false
}

func (c *Controller) ConectionAllocation(source, destination int, bitRate connections.BitRate, network infrastructure.Network, path connections.Routes, numberOfBands int) (bool, connections.Connection) {
	return c.Allocator(source, destination, bitRate, network, path, numberOfBands)
}

func (c *Controller) SetAllocator(allocator allocator.Allocator) {
	c.Allocator = allocator
}

func (c *Controller) ControllerInit(pathToRoutes string, network infrastructure.Network, allocator allocator.Allocator) {

	var connections []connections.Connection

	routes, err := c.Routes.ReadRoutesFile(pathToRoutes)

	if err != nil {
		log.Fatalf("Error reading routes file: %v", err)
	}

	c.SetRoutes(routes)
	c.SetConnections(connections)
	c.SetNetwork(network)
	c.SetAllocator(allocator)

}
