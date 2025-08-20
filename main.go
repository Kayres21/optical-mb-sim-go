package main

import (
	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/simulator"
)

func main() {

	networkPath := "files/networks/UKNet_BDM.json"
	capacitiesPath := "files/capacities/capacities.json"
	bitRatePath := "files/bitrate/bitrate.json"

	lambda := 10
	mu := 2

	numberOfBands := 1

	goalConnections := float64(100)

	var simulator simulator.Simulator

	simulator.SimulatorInit(networkPath, capacitiesPath, bitRatePath, lambda, mu, goalConnections, allocator.FirstFit, numberOfBands)
	simulator.SimulatorStart()
}
