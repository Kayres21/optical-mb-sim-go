package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/loader"
	"github.com/Kayres21/optical-mb-sim-go/internal/simulator"
)

func main() {
	networkPath := flag.String("network", "files/networks/UKNet_BDM.json", "Path to network JSON file")
	routesPath := flag.String("routes", "files/routes/UKNet_routes.json", "Path to routes JSON file")
	capacitiesPath := flag.String("capacities", "files/capacities/capacities.json", "Path to capacities JSON file")
	bitRatePath := flag.String("bitrate", "files/bitrate/bitrate.json", "Path to bitrate JSON file")
	lambda := flag.Float64("lambda", 50, "Arrival rate λ")
	mu := flag.Float64("mu", 1, "Service rate μ")
	numberOfBands := flag.Int("bands", 1, "Number of frequency bands (1–4)")
	goalConns := flag.Float64("goal", 1e8, "Number of connection requests to simulate")
	logsOn := flag.Bool("logs", true, "Enable progress logging")
	legacyOn := flag.Bool("legacy", false, "Use legacy file formats")

	flag.Parse()

	networkSet := false
	routesSet := false
	capacitiesSet := false
	bitrateSet := false

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "network":
			networkSet = true
		case "routes":
			routesSet = true
		case "capacities":
			capacitiesSet = true
		case "bitrate":
			bitrateSet = true
		}
	})

	if *legacyOn {
		if !networkSet {
			*networkPath = "legacy_files/networks/UKNet.json"
		}
		if !routesSet {
			*routesPath = "legacy_files/routes/UKNet_routes.json"
		}
		if !bitrateSet {
			*bitRatePath = "legacy_files/bitrates/bitrate_iroBand_C.json"
		}
		if !capacitiesSet {
			*capacitiesPath = ""
		}
	}

	var resLoader loader.ResourceLoader
	if *legacyOn {
		resLoader = &loader.LegacyLoader{}
	} else {
		resLoader = &loader.StandardLoader{}
	}

	network, err := resLoader.LoadNetwork(*networkPath, *capacitiesPath)
	if err != nil {
		log.Fatalf("Failed to load network: %v", err)
	}

	bitRate, err := resLoader.LoadBitRate(*bitRatePath)
	if err != nil {
		log.Fatalf("Failed to load bitrate: %v", err)
	}

	routes, err := resLoader.LoadRoutes(*routesPath)
	if err != nil {
		log.Fatalf("Failed to load routes: %v", err)
	}

	sim, err := simulator.New(
		network, bitRate, routes,
		*lambda, *mu, *goalConns,
		allocator.FirstFit,
		*numberOfBands,
	)
	if err != nil {
		log.Fatalf("Failed to initialise simulator: %v", err)
	}

	sim.Start(*logsOn)

	title := fmt.Sprintf("FirstFit_%s-erlang-%s_%s",
		network.Alias,
		fmt.Sprintf("%.1f", *lambda),
		strconv.Itoa(*numberOfBands),
	)
	if err := sim.Plot(title, "Número de conexiones", "Probabilidad de bloqueo"); err != nil {
		log.Fatalf("Failed to generate plot: %v", err)
	}
}

// 1 banda lambda 50 mu 1: 0.690987 00:04:02 1e8
// 2 banda lambda 50 mu 1: 0.518572 00:07:09 1e8
// 3 banda lambda 50 mu 1: 0.387643 00:10:18 1e8
// 4 banda lambda 50 mu 1: 0.289645 00:15:34 1e8
