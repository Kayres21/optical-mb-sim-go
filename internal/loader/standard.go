package loader

import (
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

// StandardLoader implements ResourceLoader for the current standard file formats.
type StandardLoader struct{}

func (l *StandardLoader) LoadNetwork(networkPath, capacitiesPath string) (infrastructure.Network, error) {
	return infrastructure.NetworkGenerate(networkPath, capacitiesPath)
}

func (l *StandardLoader) LoadBitRate(bitRatePath string) (connections.BitRateList, error) {
	return connections.ReadBitRateFile(bitRatePath)
}

func (l *StandardLoader) LoadRoutes(routesPath string) (connections.Routes, error) {
	return connections.ReadRoutesFile(routesPath)
}
