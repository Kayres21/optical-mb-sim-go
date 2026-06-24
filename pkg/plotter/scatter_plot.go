package plotter

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"time"

	gplot "gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// PlotConfig holds optional styling overrides for GenerateLinePlot.
type PlotConfig struct {
	// LineColor is the colour of the data line. Defaults to a vivid blue.
	LineColor color.RGBA
	// PointColor is the colour of the scatter markers. Defaults to the same as LineColor.
	PointColor color.RGBA
	// Width and Height of the output image in inches.
	Width, Height vg.Length
	// OutputDir is the directory where the PNG is saved. Defaults to "result".
	OutputDir string
}

var defaultConfig = PlotConfig{
	LineColor:  color.RGBA{R: 30, G: 120, B: 255, A: 255},
	PointColor: color.RGBA{R: 30, G: 120, B: 255, A: 255},
	Width:      8 * vg.Inch,
	Height:     5 * vg.Inch,
	OutputDir:  "result",
}

// GenerateScatterPlot is the original API kept for backwards compatibility.
// It now delegates to GenerateLinePlot with default styling.
func GenerateScatterPlot(xData, yData []float64, title, xLabel, yLabel string) error {
	return GenerateLinePlot(xData, yData, title, xLabel, yLabel, defaultConfig)
}

// GenerateLinePlot produces a polished line+scatter plot and saves it as PNG.
// The filename is "<title>_<timestamp>.png" inside cfg.OutputDir.
func GenerateLinePlot(xData, yData []float64, title, xLabel, yLabel string, cfg PlotConfig) error {
	if len(xData) != len(yData) {
		return fmt.Errorf("x and y slices must have the same length (got %d vs %d)", len(xData), len(yData))
	}
	if len(xData) == 0 {
		return fmt.Errorf("no data points to plot")
	}

	// ── Build XYs ────────────────────────────────────────────────────────────
	pts := make(plotter.XYs, len(xData))
	for i := range xData {
		pts[i].X = xData[i]
		pts[i].Y = yData[i]
	}

	// ── Create plot ───────────────────────────────────────────────────────────
	p := gplot.New()

	p.Title.Text = title
	p.Title.Padding = vg.Points(10)

	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel
	p.X.Label.Padding = vg.Points(6)
	p.Y.Label.Padding = vg.Points(6)

	// Axis tick formatters
	p.X.Tick.Marker = scientificTicker{}
	p.Y.Tick.Marker = gplot.DefaultTicks{}

	// Grid
	grid := plotter.NewGrid()
	grid.Horizontal.Color = color.RGBA{R: 220, G: 220, B: 220, A: 255}
	grid.Vertical.Color = color.RGBA{R: 220, G: 220, B: 220, A: 255}
	grid.Horizontal.Width = vg.Points(0.5)
	grid.Vertical.Width = vg.Points(0.5)
	p.Add(grid)

	// ── Line ──────────────────────────────────────────────────────────────────
	line, err := plotter.NewLine(pts)
	if err != nil {
		return fmt.Errorf("creating line: %w", err)
	}
	line.Color = cfg.LineColor
	line.Width = vg.Points(2)
	p.Add(line)

	// ── Scatter markers ───────────────────────────────────────────────────────
	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		return fmt.Errorf("creating scatter: %w", err)
	}
	scatter.GlyphStyle = draw.GlyphStyle{
		Color:  cfg.PointColor,
		Radius: vg.Points(4),
		Shape:  draw.CircleGlyph{},
	}
	p.Add(scatter)

	// ── Y-axis range: pad slightly beyond data bounds ─────────────────────────
	minY, maxY := minMax(yData)
	padding := (maxY - minY) * 0.1
	if padding == 0 {
		padding = 0.05
	}
	p.Y.Min = math.Max(0, minY-padding)
	p.Y.Max = maxY + padding

	// ── Save ──────────────────────────────────────────────────────────────────
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.png", title, timestamp)
	filePath := filepath.Join(cfg.OutputDir, filename)

	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory %q: %w", cfg.OutputDir, err)
	}

	if err := p.Save(cfg.Width, cfg.Height, filePath); err != nil {
		return fmt.Errorf("saving plot to %q: %w", filePath, err)
	}

	fmt.Printf("Plot saved → %s\n", filePath)
	return nil
}

// GenerateHistogram produces a polished histogram plot and saves it as PNG.
func GenerateHistogram(data []float64, bins int, title, xLabel, yLabel string) error {
	if len(data) == 0 {
		return fmt.Errorf("no data points to plot")
	}

	p := gplot.New()
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	// Create a histogram of our values.
	v := make(plotter.Values, len(data))
	copy(v, data)

	h, err := plotter.NewHist(v, bins)
	if err != nil {
		return fmt.Errorf("creating histogram: %w", err)
	}
	h.Color = color.RGBA{R: 30, G: 120, B: 255, A: 255}
	p.Add(h)

	// Save the plot.
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.png", title, timestamp)
	filePath := filepath.Join(defaultConfig.OutputDir, filename)

	if err := os.MkdirAll(defaultConfig.OutputDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory %q: %w", defaultConfig.OutputDir, err)
	}

	if err := p.Save(defaultConfig.Width, defaultConfig.Height, filePath); err != nil {
		return fmt.Errorf("saving plot to %q: %w", filePath, err)
	}

	fmt.Printf("Histogram saved → %s\n", filePath)
	return nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

func minMax(data []float64) (min, max float64) {
	min, max = data[0], data[0]
	for _, v := range data[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
}

// scientificTicker formats large X values in scientific notation (e.g. 1×10⁷).
type scientificTicker struct{}

func (scientificTicker) Ticks(min, max float64) []gplot.Tick {
	ticks := gplot.DefaultTicks{}.Ticks(min, max)
	for i, t := range ticks {
		if t.Label == "" {
			continue
		}
		v := t.Value
		if math.Abs(v) >= 1e6 {
			exp := int(math.Log10(math.Abs(v)))
			mantissa := v / math.Pow10(exp)
			ticks[i].Label = fmt.Sprintf("%.0f×10⁷", mantissa*10/math.Pow10(exp-6))
			// simpler: just show e-notation
			ticks[i].Label = fmt.Sprintf("%.2g", v)
		}
	}
	return ticks
}
