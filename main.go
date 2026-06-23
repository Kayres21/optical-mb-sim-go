package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/defragmentator"
	"github.com/Kayres21/optical-mb-sim-go/internal/loader"
	"github.com/Kayres21/optical-mb-sim-go/internal/simulator"
)

type AppConfig struct {
	Network    string   `json:"network"`
	Routes     string   `json:"routes"`
	Capacities string   `json:"capacities"`
	Bitrate    string   `json:"bitrate"`
	Lambda     *float64 `json:"lambda"`
	Mu         *float64 `json:"mu"`
	Bands      *int     `json:"bands"`
	Goal       *float64 `json:"goal"`
	Logs       *bool    `json:"logs"`
	Legacy     *bool    `json:"legacy"`
	DefragMode string   `json:"defrag_mode"`
}

func loadConfig(path string) (AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return AppConfig{}, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}

func defaultConfig() AppConfig {
	lambda := 50.0
	mu := 1.0
	bands := 1
	goal := 1e8
	logs := true
	legacy := false

	return AppConfig{
		Network:    "files/networks/UKNet_BDM.json",
		Routes:     "files/routes/UKNet_routes.json",
		Capacities: "files/capacities/capacities.json",
		Bitrate:    "files/bitrate/bitrate.json",
		Lambda:     &lambda,
		Mu:         &mu,
		Bands:      &bands,
		Goal:       &goal,
		Logs:       &logs,
		Legacy:     &legacy,
		DefragMode: defragmentator.DefragNone,
	}
}

func applyDefaults(cfg AppConfig) AppConfig {
	defaults := defaultConfig()

	if cfg.Network == "" {
		cfg.Network = defaults.Network
	}
	if cfg.Routes == "" {
		cfg.Routes = defaults.Routes
	}
	if cfg.Capacities == "" {
		cfg.Capacities = defaults.Capacities
	}
	if cfg.Bitrate == "" {
		cfg.Bitrate = defaults.Bitrate
	}
	if cfg.Lambda == nil {
		cfg.Lambda = defaults.Lambda
	}
	if cfg.Mu == nil {
		cfg.Mu = defaults.Mu
	}
	if cfg.Bands == nil {
		cfg.Bands = defaults.Bands
	}
	if cfg.Goal == nil {
		cfg.Goal = defaults.Goal
	}
	if cfg.Logs == nil {
		cfg.Logs = defaults.Logs
	}
	if cfg.Legacy == nil {
		cfg.Legacy = defaults.Legacy
	}
	if cfg.DefragMode == "" {
		cfg.DefragMode = defaults.DefragMode
	}

	return cfg
}

func main() {
	configPath := flag.String("config", "files/config.json", "Path to JSON configuration file")
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cfg = applyDefaults(cfg)

	var resLoader loader.ResourceLoader
	if *cfg.Legacy {
		resLoader = &loader.LegacyLoader{}
	} else {
		resLoader = &loader.StandardLoader{}
	}

	network, err := resLoader.LoadNetwork(cfg.Network, cfg.Capacities)
	if err != nil {
		log.Fatalf("Failed to load network: %v", err)
	}

	bitRate, err := resLoader.LoadBitRate(cfg.Bitrate)
	if err != nil {
		log.Fatalf("Failed to load bitrate: %v", err)
	}

	routes, err := resLoader.LoadRoutes(cfg.Routes)
	if err != nil {
		log.Fatalf("Failed to load routes: %v", err)
	}

	sim, err := simulator.New(
		network, bitRate, routes,
		*cfg.Lambda, *cfg.Mu, *cfg.Goal,
		allocator.FirstFit,
		*cfg.Bands,
		cfg.DefragMode,
		defragmentator.DefaultDecision,
		defragmentator.DefaultAction,
	)
	if err != nil {
		log.Fatalf("Failed to initialise simulator: %v", err)
	}

	sim.Start(*cfg.Logs)

	title := fmt.Sprintf("FirstFit_%s-erlang-%s_%s",
		network.Alias,
		fmt.Sprintf("%.1f", *cfg.Lambda),
		strconv.Itoa(*cfg.Bands),
	)
	if err := sim.Plot(title, "Número de conexiones", "Probabilidad de bloqueo"); err != nil {
		log.Fatalf("Failed to generate plot: %v", err)
	}
}

// 1 banda lambda 50 mu 1: 0.690987 00:04:02 1e8
// 2 banda lambda 50 mu 1: 0.518572 00:07:09 1e8
// 3 banda lambda 50 mu 1: 0.387643 00:10:18 1e8
// 4 banda lambda 50 mu 1: 0.289645 00:15:34 1e8
