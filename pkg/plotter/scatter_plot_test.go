package plotter

import (
	"strings"
	"testing"
)

func TestGenerateMultiSeriesPlotRejectsMismatchedLengths(t *testing.T) {
	series := []Series{
		{Label: "A", X: []float64{1, 2, 3}, Y: []float64{0.1, 0.2, 0.3}},
		{Label: "B", X: []float64{1, 2}, Y: []float64{0.2, 0.4, 0.6}},
	}

	err := GenerateMultiSeriesPlot(series, "test", "x", "y", defaultConfig)
	if err == nil || !strings.Contains(err.Error(), "mismatched") {
		t.Fatalf("GenerateMultiSeriesPlot() error = %v, want length mismatch error", err)
	}
}
