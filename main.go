package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/defragmentator"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
	"github.com/Kayres21/optical-mb-sim-go/internal/loader"
	"github.com/Kayres21/optical-mb-sim-go/internal/simulator"
	"github.com/Kayres21/optical-mb-sim-go/simulations"
)

type AppConfig struct {
	Network         string   `json:"network"`
	Routes          string   `json:"routes"`
	Capacities      string   `json:"capacities"`
	Bitrate         string   `json:"bitrate"`
	Lambda          *float64 `json:"lambda"`
	Mu              *float64 `json:"mu"`
	Bands           *int     `json:"bands"`
	Goal            *float64 `json:"goal"`
	Logs            *bool    `json:"logs"`
	Legacy          *bool    `json:"legacy"`
	DefragMode      string   `json:"defrag_mode"`
	DefragAlgorithm string   `json:"defrag_algorithm"`
	EventsCSV       string   `json:"events_csv"`
	Sweep           *bool    `json:"sweep"`
	LambdaStart     *float64 `json:"lambda_start"`
	LambdaEnd       *float64 `json:"lambda_end"`
	LambdaStep      *float64 `json:"lambda_step"`
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
		Network:         "files/networks/UKNet_BDM.json",
		Routes:          "files/routes/UKNet_routes.json",
		Capacities:      "files/capacities/capacities.json",
		Bitrate:         "files/bitrate/bitrate.json",
		Lambda:          &lambda,
		Mu:              &mu,
		Bands:           &bands,
		Goal:            &goal,
		Logs:            &logs,
		Legacy:          &legacy,
		DefragMode:      defragmentator.DefragBeforeArrival,
		DefragAlgorithm: "first_fit",
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
	if cfg.DefragAlgorithm == "" {
		cfg.DefragAlgorithm = defaults.DefragAlgorithm
	}
	if cfg.EventsCSV == "" {
		cfg.EventsCSV = ""
	}
	if cfg.Sweep == nil {
		cfg.Sweep = new(bool)
		*cfg.Sweep = false
	}
	if cfg.LambdaStart == nil {
		cfg.LambdaStart = new(float64)
		*cfg.LambdaStart = 500
	}
	if cfg.LambdaEnd == nil {
		cfg.LambdaEnd = new(float64)
		*cfg.LambdaEnd = 1500
	}
	if cfg.LambdaStep == nil {
		cfg.LambdaStep = new(float64)
		*cfg.LambdaStep = 50
	}

	return cfg
}

// setupFileLogging creates a timestamped log file under logs/ and tees all
// stdout, stderr and slog output to it, returning a cleanup func to be
// deferred.
func setupFileLogging() (func(), error) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logPath := filepath.Join("logs", fmt.Sprintf("simulation_%s.log", time.Now().Format("20060102_150405")))
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	origStdout := os.Stdout
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		logFile.Close()
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	os.Stdout = stdoutW

	origStderr := os.Stderr
	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		stdoutW.Close()
		stdoutR.Close()
		logFile.Close()
		os.Stdout = origStdout
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	os.Stderr = stderrW

	stdoutDone := make(chan struct{})
	go func() {
		defer close(stdoutDone)
		io.Copy(io.MultiWriter(origStdout, logFile), stdoutR)
	}()
	stderrDone := make(chan struct{})
	go func() {
		defer close(stderrDone)
		io.Copy(io.MultiWriter(origStderr, logFile), stderrR)
	}()

	log.SetOutput(io.MultiWriter(origStderr, logFile))
	slog.SetDefault(slog.New(slog.NewTextHandler(io.MultiWriter(origStderr, logFile), nil)))

	cleanup := func() {
		stdoutW.Close()
		stderrW.Close()
		<-stdoutDone
		<-stderrDone
		logFile.Close()
		os.Stdout = origStdout
		os.Stderr = origStderr
	}
	return cleanup, nil
}

func main() {
	cleanupLogging, err := setupFileLogging()
	if err != nil {
		log.Fatalf("Failed to set up file logging: %v", err)
	}
	defer cleanupLogging()

	configPath := flag.String("config", "files/config.json", "Path to JSON configuration file")
	eventsCSV := flag.String("events-csv", "", "Path to write the generated events CSV after the simulation")
	logs := flag.Bool("logs", true, "Enable progress logging")
	defragMode := flag.String("defrag-mode", "", "Defragmentation mode: none, before_arrival, after_block, after_assign")
	lambdaStart := flag.Float64("lambda-start", 500, "Start lambda for a multi-simulation sweep")
	lambdaEnd := flag.Float64("lambda-end", 1500, "End lambda for a multi-simulation sweep")
	lambdaStep := flag.Float64("lambda-step", 50, "Lambda step for a multi-simulation sweep")
	sweepEnabled := flag.Bool("sweep", false, "Run multiple simulations across a lambda range and plot all results together")
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cfg = applyDefaults(cfg)
	if *defragMode != "" {
		cfg.DefragMode = *defragMode
	}
	if *eventsCSV != "" {
		cfg.EventsCSV = *eventsCSV
	}
	if cfg.Logs == nil || *cfg.Logs != *logs {
		cfg.Logs = logs
	}

	allowedModes := map[string]bool{
		defragmentator.DefragNone:          true,
		defragmentator.DefragBeforeArrival: true,
		defragmentator.DefragAfterBlock:    true,
		defragmentator.DefragAfterAssign:   true,
	}
	if !allowedModes[cfg.DefragMode] {
		log.Fatalf("Unsupported defrag mode %q. Valid modes: none, before_arrival, after_block, after_assign", cfg.DefragMode)
	}
	allowedAlgorithms := map[string]bool{
		"first_fit":           true,
		"multiband":           true,
		"multiband_same_band": true,
	}
	if !allowedAlgorithms[cfg.DefragAlgorithm] {
		log.Fatalf("Unsupported defrag algorithm %q. Valid algorithms: first_fit, multiband, multiband_same_band", cfg.DefragAlgorithm)
	}

	defragDecision := defragmentator.DefaultDecision
	defragAction := defragmentator.DefaultAction
	if cfg.DefragMode != defragmentator.DefragNone {
		switch cfg.DefragAlgorithm {
		case "multiband":
			defragAction = func(network infrastructure.Network, connectionsMap map[string]connections.Connection, routes connections.Routes, alloc allocator.Allocator, numberOfBands int) (int, error) {
				return defragmentator.MultiBandActiveConnections(network, connectionsMap, routes, numberOfBands)
			}
		case "multiband_same_band":
			defragAction = func(network infrastructure.Network, connectionsMap map[string]connections.Connection, routes connections.Routes, alloc allocator.Allocator, numberOfBands int) (int, error) {
				return defragmentator.MultiBandSameBandActiveConnections(network, connectionsMap, routes, numberOfBands)
			}
		default:
			defragAction = func(network infrastructure.Network, connectionsMap map[string]connections.Connection, routes connections.Routes, alloc allocator.Allocator, numberOfBands int) (int, error) {
				return defragmentator.FirstFitActiveConnections(network, connectionsMap, routes, numberOfBands)
			}
		}
	}

	fmt.Printf("Running simulation with defrag mode: %s, algorithm: %s\n", cfg.DefragMode, cfg.DefragAlgorithm)

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

	bitRate, err := resLoader.LoadBitRate(cfg.Bitrate, *cfg.Bands)
	if err != nil {
		log.Fatalf("Failed to load bitrate: %v", err)
	}

	routes, err := resLoader.LoadRoutes(cfg.Routes)
	if err != nil {
		log.Fatalf("Failed to load routes: %v", err)
	}

	if *sweepEnabled || (cfg.Sweep != nil && *cfg.Sweep) {
		if *sweepEnabled {
			cfg.Sweep = sweepEnabled
		}
		if cfg.LambdaStart == nil || cfg.LambdaEnd == nil || cfg.LambdaStep == nil {
			cfg.LambdaStart = lambdaStart
			cfg.LambdaEnd = lambdaEnd
			cfg.LambdaStep = lambdaStep
		}

		runner := simulations.NewSweepRunner(network, bitRate, routes, simulations.SweepConfig{
			StartLambda: *cfg.LambdaStart,
			EndLambda:   *cfg.LambdaEnd,
			StepLambda:  *cfg.LambdaStep,
			Mu:          *cfg.Mu,
			Bands:       *cfg.Bands,
			Goal:        *cfg.Goal,
			LogOn:       *cfg.Logs,
		})
		if _, err := runner.Run(); err != nil {
			log.Fatalf("Failed to run lambda sweep: %v", err)
		}

		title := fmt.Sprintf("FirstFit_%s_lambda_sweep_%s_%s", network.Alias,
			strconv.FormatInt(int64(*cfg.LambdaStart), 10),
			strconv.FormatInt(int64(*cfg.LambdaEnd), 10),
		)
		if err := runner.Plot(title, "Número de conexiones", "Probabilidad de bloqueo"); err != nil {
			log.Fatalf("Failed to generate combined plot: %v", err)
		}
		return
	}

	sim, err := simulator.New(
		network, bitRate, routes,
		*cfg.Lambda, *cfg.Mu, *cfg.Goal,
		allocator.FirstFit,
		*cfg.Bands,
		cfg.DefragMode,
		defragDecision,
		defragAction,
	)
	if err != nil {
		log.Fatalf("Failed to initialise simulator: %v", err)
	}
	sim.SetRecordEvents(cfg.EventsCSV != "")
	sim.Start(*cfg.Logs)

	if cfg.EventsCSV != "" {
		if err := sim.SaveEventsCSV(cfg.EventsCSV); err != nil {
			log.Fatalf("Failed to save events CSV: %v", err)
		}
	}

	title := fmt.Sprintf("FirstFit_%s-erlang-%s_%s",
		network.Alias,
		fmt.Sprintf("%.1f", *cfg.Lambda),
		strconv.Itoa(*cfg.Bands),
	)
	if err := sim.Plot(title, "Número de conexiones", "Probabilidad de bloqueo"); err != nil {
		log.Fatalf("Failed to generate plot: %v", err)
	}
}
