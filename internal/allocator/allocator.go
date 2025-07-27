package allocator

import (
	"log"
	"simulator/internal/connections"
	"simulator/internal/infrastructure"
)

type Allocator func(source, destination string, bitRate connections.BitRate, network infrastructure.Network, path connections.Routes, numberOfBands int) (bool, connections.Connection)

func FirstFit(source int, destination int, bitRate connections.BitRate, network infrastructure.Network, path connections.Routes, numberOfBands int) (bool, connections.Connection) {
	pathSelected := path.GetKshortestPath(0, source, destination)
	// Implement the allocation logic here
	// This is a placeholder implementation
	log.Printf("Allocating connection from %d to %d with bit rate %v and path %d", source, destination, bitRate, pathSelected)

	connection := connections.Connection{}

	return true, connection // Return true if allocation is successful
}
