package loader

import (
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

// ResourceLoader defines the interface for loading simulation resources.
type ResourceLoader interface {
	LoadNetwork(networkPath, capacitiesPath string) (infrastructure.Network, error)
	// LoadBitRate loads the bitrate file. numberOfBands selects which
	// per-band-count reach configuration to use for schema-based files.
	LoadBitRate(bitRatePath string, numberOfBands int) (connections.BitRateList, error)
	LoadRoutes(routesPath string) (connections.Routes, error)
}
