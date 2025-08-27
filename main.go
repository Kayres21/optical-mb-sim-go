package main

import (
	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/simulator"
	"github.com/Kayres21/optical-mb-sim-go/pkg/plotter"
	"github.com/google/uuid"
)

func main() {

	networkPath := "files/networks/UKNet_BDM.json"
	capacitiesPath := "files/capacities/capacities.json"
	bitRatePath := "files/bitrate/bitrate.json"
	routesPath := "files/routes/UKNet_routes.json"

	lambda := 5
	mu := 1

	numberOfBands := 1

	goalConnections := float64(1e6)

	var simulator simulator.Simulator

	simulator.SimulatorInit(networkPath, routesPath, capacitiesPath, bitRatePath, lambda, mu, goalConnections, allocator.FirstFit, numberOfBands)
	simulator.SimulatorStart(true)
	plotter.GenerateScatterPlot(simulator.GetArrives(), simulator.GetResults(), "FirstFit_UKNet-"+uuid.NewString(), "Número de conexiones", "Probabilidad de bloqueo")
}
