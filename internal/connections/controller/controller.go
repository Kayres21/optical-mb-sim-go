package controller

import (
	"fmt"
	"log"

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

func (c *Controller) SetAllocator(allocator allocator.Allocator) {
	c.Allocator = allocator
}

func (c *Controller) RoutesInit(pathToRoutes string) {

	routes, err := c.Routes.ReadRoutesFile(pathToRoutes)

	if err != nil {
		log.Fatalf("Error reading routes file: %v", err)
	}

	c.SetRoutes(routes)
}

func (c *Controller) ConnectionsInit() {
	var connections []connections.Connection
	c.SetConnections(connections)
}

func (c *Controller) ControllerInit(pathToRoutes string, network infrastructure.Network, allocator allocator.Allocator) {
	c.RoutesInit(pathToRoutes)
	c.ConnectionsInit()
	c.SetNetwork(network)
	c.SetAllocator(allocator)

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

	links := con.GetLinks()

	for _, link := range links {
		link.ReleaseConnection(con.GetInitialSlot(), con.GetSlots(), con.GetBandSelected())

	}

	return nil
}

func (c *Controller) ConectionAllocation(source, destination int, slot int, network infrastructure.Network, path connections.Routes, numberOfBands int, id string) (bool, connections.Connection) {
	return c.Allocator(source, destination, slot, network, path, numberOfBands, id)
}
