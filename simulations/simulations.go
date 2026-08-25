package simulations

import (
	"fmt"
	"math"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/defragmentator"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
	"github.com/Kayres21/optical-mb-sim-go/internal/simulator"
	"github.com/Kayres21/optical-mb-sim-go/pkg/plotter"
)

type RunResult struct {
	Lambda  float64
	Arrives []float64
	Results []float64
}

type SweepConfig struct {
	StartLambda float64
	EndLambda   float64
	StepLambda  float64
	Mu          float64
	Bands       int
	Goal        float64
	LogOn       bool
}

type SweepRunner struct {
	Network    infrastructure.Network
	BitRate    connections.BitRateList
	Routes     connections.Routes
	Config     SweepConfig
	RunResults []RunResult
	PlotConfig plotter.PlotConfig
}

func NewSweepRunner(
	network infrastructure.Network,
	bitRate connections.BitRateList,
	routes connections.Routes,
	config SweepConfig,
) *SweepRunner {
	return &SweepRunner{
		Network:    network,
		BitRate:    bitRate,
		Routes:     routes,
		Config:     config,
		PlotConfig: plotter.DefaultPlotConfig(),
	}
}

func (r *SweepRunner) Run() ([]RunResult, error) {
	cfg := r.Config
	if cfg.StepLambda <= 0 {
		return nil, fmt.Errorf("lambda step must be > 0")
	}
	if cfg.StartLambda > cfg.EndLambda {
		return nil, fmt.Errorf("start lambda must be <= end lambda")
	}

	count := int(math.Floor((cfg.EndLambda-cfg.StartLambda)/cfg.StepLambda)) + 1
	results := make([]RunResult, 0, count)

	for lambda := cfg.StartLambda; lambda <= cfg.EndLambda+1e-9; lambda += cfg.StepLambda {
		sim, err := simulator.New(
			r.Network,
			r.BitRate,
			r.Routes,
			lambda,
			cfg.Mu,
			cfg.Goal,
			allocator.FirstFit,
			cfg.Bands,
			defragmentator.DefragNone,
			defragmentator.DefaultDecision,
			defragmentator.DefaultAction,
		)
		if err != nil {
			return nil, fmt.Errorf("initialising simulator for lambda %.0f: %w", lambda, err)
		}

		sim.Start(cfg.LogOn)
		results = append(results, RunResult{
			Lambda:  lambda,
			Arrives: sim.Arrives(),
			Results: sim.Results(),
		})
	}

	r.RunResults = results
	return results, nil
}

func (r *SweepRunner) BuildSeries() []plotter.Series {
	series := make([]plotter.Series, len(r.RunResults))
	for i, res := range r.RunResults {
		series[i] = plotter.Series{
			Label: fmt.Sprintf("lambda=%.0f", res.Lambda),
			X:     res.Arrives,
			Y:     res.Results,
		}
	}
	return series
}

func (r *SweepRunner) Plot(title, xLabel, yLabel string) error {
	if len(r.RunResults) == 0 {
		return fmt.Errorf("no sweep results to plot")
	}
	r.PlotConfig.OutputDir = "result"
	return plotter.GenerateMultiSeriesPlot(r.BuildSeries(), title, xLabel, yLabel, r.PlotConfig)
}

func RunLambdaSweep(
	network infrastructure.Network,
	bitRate connections.BitRateList,
	routes connections.Routes,
	startLambda, endLambda, stepLambda, mu float64,
	bands int,
	goalConnections float64,
	logOn bool,
) ([]RunResult, error) {
	runner := NewSweepRunner(network, bitRate, routes, SweepConfig{
		StartLambda: startLambda,
		EndLambda:   endLambda,
		StepLambda:  stepLambda,
		Mu:          mu,
		Bands:       bands,
		Goal:        goalConnections,
		LogOn:       logOn,
	})
	return runner.Run()
}
