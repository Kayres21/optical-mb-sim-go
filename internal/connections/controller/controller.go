package controller

import (
	"log"
	"simulator/internal/connections"
	"simulator/internal/infrastructure"
)

type Controller struct {
	Routes      connections.Routes
	Connections []connections.Connection
	Network     infrastructure.Network
	Allocator   Allocator
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

type Allocator func(source, destination string, bitRate connections.BitRate, connection connections.Connection, network infrastructure.Network, path connections.Routes) (bool, connections.Connection)

func (c *Controller) ConectionAllocation(source, destination string, bitRate connections.BitRate, connection connections.Connection, network infrastructure.Network, path connections.Routes) (bool, connections.Connection) {
	return c.Allocator(source, destination, bitRate, connection, network, path)
}

func (c *Controller) ControllerInit(pathToRoutes string, network infrastructure.Network, connections []connections.Connection, allocator Allocator) {

	routes, err := c.Routes.ReadRoutesFile(pathToRoutes)

	if err != nil {
		log.Fatalf("Error reading routes file: %v", err)
	}

	c.SetRoutes(routes)
	c.SetConnections(connections)
	c.SetNetwork(network)
	c.Allocator = allocator

}
