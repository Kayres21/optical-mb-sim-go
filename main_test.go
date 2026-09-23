package main

import (
	"flag"
	"strings"
	"testing"
)

func TestSimulationLogNameIncludesLambdaAndMu(t *testing.T) {
	got := simulationLogName("EON", 4, "multiband", 50.5, 1.25)

	if !strings.Contains(got, "lambda_50.5") {
		t.Fatalf("log name missing lambda: %q", got)
	}
	if !strings.Contains(got, "mu_1.25") {
		t.Fatalf("log name missing mu: %q", got)
	}
}

func TestApplyExplicitSweepFlagsOverrideConfigValues(t *testing.T) {
	cfg := AppConfig{
		Sweep:       boolPtr(true),
		LambdaStart: floatPtr(500),
		LambdaEnd:   floatPtr(2000),
		LambdaStep:  floatPtr(100),
	}

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	sweepEnabled := fs.Bool("sweep", false, "")
	lambdaStart := fs.Float64("lambda-start", 500, "")
	lambdaEnd := fs.Float64("lambda-end", 1500, "")
	lambdaStep := fs.Float64("lambda-step", 50, "")

	if err := fs.Parse([]string{"-sweep=true", "-lambda-start=4000", "-lambda-end=4000", "-lambda-step=10"}); err != nil {
		t.Fatalf("fs.Parse() error = %v", err)
	}

	cfg = applyExplicitSweepFlags(cfg, fs, sweepEnabled, lambdaStart, lambdaEnd, lambdaStep)

	if *cfg.Sweep != true {
		t.Fatalf("expected sweep override to true, got %v", *cfg.Sweep)
	}
	if *cfg.LambdaStart != 4000 {
		t.Fatalf("expected lambda start override to 4000, got %v", *cfg.LambdaStart)
	}
	if *cfg.LambdaEnd != 4000 {
		t.Fatalf("expected lambda end override to 4000, got %v", *cfg.LambdaEnd)
	}
	if *cfg.LambdaStep != 10 {
		t.Fatalf("expected lambda step override to 10, got %v", *cfg.LambdaStep)
	}
}

func boolPtr(v bool) *bool {
	return &v
}

func floatPtr(v float64) *float64 {
	return &v
}
