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

	goalConnections := 100000000000

	var simulator simulator.Simulator

	simulator.SimulatorInit(networkPath, capacitiesPath, bitRatePath, lambda, mu, goalConnections)

}
