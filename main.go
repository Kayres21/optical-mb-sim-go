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

	lambda := 500
	mu := 1

	numberOfBands := 4

	goalConnections := 1e6

	logsOn := true

	var simulator simulator.Simulator

	simulator.SimulatorInit(networkPath, routesPath, capacitiesPath, bitRatePath, lambda, mu, goalConnections, allocator.FirstFit, numberOfBands)
	simulator.SimulatorStart(logsOn)
	simulator.SimulatorPlot("FirstFit_UKNet-erlang-"+strconv.Itoa(lambda), "Número de conexiones", "Probabilidad de bloqueo")
}
