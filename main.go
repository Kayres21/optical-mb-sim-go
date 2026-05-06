package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/simulator"
)

func main() {
	networkPath := flag.String("network", "files/networks/UKNet_BDM.json", "Path to network JSON file")
	routesPath := flag.String("routes", "files/routes/UKNet_routes.json", "Path to routes JSON file")
	capacitiesPath := flag.String("capacities", "files/capacities/capacities.json", "Path to capacities JSON file")
	bitRatePath := flag.String("bitrate", "files/bitrate/bitrate.json", "Path to bitrate JSON file")
	lambda := flag.Int("lambda", 50, "Arrival rate λ")
	mu := flag.Int("mu", 1, "Service rate μ")
	numberOfBands := flag.Int("bands", 1, "Number of frequency bands (1–4)")
	goalConns := flag.Float64("goal", 1e8, "Number of connection requests to simulate")
	logsOn := flag.Bool("logs", true, "Enable progress logging")
	flag.Parse()

	sim, err := simulator.New(
		*networkPath, *routesPath, *capacitiesPath, *bitRatePath,
		*lambda, *mu, *goalConns,
		allocator.FirstFit,
		*numberOfBands,
	)
	if err != nil {
		log.Fatalf("Failed to initialise simulator: %v", err)
	}

	sim.Start(*logsOn)

	title := fmt.Sprintf("FirstFit_UKNet-erlang-%s_%s",
		strconv.Itoa(*lambda),
		strconv.Itoa(*numberOfBands),
	)
	if err := sim.Plot(title, "Número de conexiones", "Probabilidad de bloqueo"); err != nil {
		log.Fatalf("Failed to generate plot: %v", err)
	}
}

// 1 banda lambda 50 mu 1: 0.690986 00:04:02 1e8
// 2 banda lambda 50 mu 1: 0.518540 00:06:22 1e8
// 3 banda lambda 50 mu 1: 0.387696 00:08:42 1e8
// 4 banda lambda 50 mu 1: 0.289625 00:09:46 1e8
