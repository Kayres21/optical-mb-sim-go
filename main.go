package main

import (
	"simulator/internal/simulator"
)

func main() {

	networkPath := "files/networks/UKNet_BDM.json"
	capacitiesPath := "files/capacities/capacities.json"
	bitRatePath := "files/bitrate/bitrate.json"

	lambda := 10
	mu := 2

	numberOfBands := 1

	goalConnections := 1e6

	var simulator simulator.Simulator

	simulator.SimulatorInit(networkPath, capacitiesPath, bitRatePath, lambda, mu, goalConnections, simulator.Controller.Allocator, numberOfBands)

}
