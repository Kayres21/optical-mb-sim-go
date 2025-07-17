package controller

import (
	"simulator/internal/connections"
	"simulator/internal/infrastructure"
)

type Controller struct {
	Routes      connections.Routes
	Connections []connections.Connection
	Network     infrastructure.Network
	Paths       connections.Routes
	Allocator   Allocator
}

type Allocator func(source, destination string, bitRate connections.BitRate, connection connections.Connection, network infrastructure.Network, path connections.Routes) (bool, connections.Connection)

func (c *Controller) ConectionAllocation(source, destination string, bitRate connections.BitRate, connection connections.Connection, network infrastructure.Network, path connections.Routes) (bool, connections.Connection) {
	return c.Allocator(source, destination, bitRate, connection, network, path)
}
