package main

import (
	"strconv"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/simulator"
)

func main() {

	networkPath := "files/networks/UKNet_BDM.json"
	capacitiesPath := "files/capacities/capacities.json"
	bitRatePath := "files/bitrate/bitrate.json"
	routesPath := "files/routes/UKNet_routes.json"

	// networkPath := "files/networks/network_test.json"
	// capacitiesPath := "files/capacities/capacities_test.json"
	// bitRatePath := "files/bitrate/bitrate_test.json"
	// routesPath := "files/routes/network_test_routes.json"

	lambda := 50
	mu := 1

	numberOfBands := 4
	goalConnections := 1e7
	logsOn := true

	sim := simulator.New(networkPath, routesPath, capacitiesPath, bitRatePath, lambda, mu, goalConnections, allocator.FirstFit, numberOfBands)
	sim.Start(logsOn)
	sim.Plot("FirstFit_UKNet-erlang-"+strconv.Itoa(lambda)+"_"+strconv.Itoa(numberOfBands), "Número de conexiones", "Probabilidad de bloqueo")
}
