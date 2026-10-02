package simulator

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"
)

const (
	fragFormatVersion = 1
	fragBufferSize    = 1 << 20
)

// fragSink streams fragmentation samples to disk as little-endian float32 rows
// ([sample][series], series = bands..., network) so long runs never hold the
// series in memory. A JSON sidecar describes the layout, and a second binary
// file records defragmentation events as (arrival, moved) uint64 pairs.
type fragSink struct {
	prefix string
	stride int

	dataFile   *os.File
	data       *bufio.Writer
	defragFile *os.File
	defrag     *bufio.Writer

	row         []byte
	pair        [16]byte
	samples     int64
	defragCount int64

	seriesNames []string
	extra       map[string]any
	startedAt   time.Time
}

type fragMeta struct {
	FormatVersion int            `json:"format_version"`
	DType         string         `json:"dtype"`
	Layout        string         `json:"layout"`
	DataFile      string         `json:"data_file"`
	NSeries       int            `json:"n_series"`
	SeriesNames   []string       `json:"series_names"`
	NSamples      int64          `json:"n_samples"`
	Stride        int            `json:"stride"`
	FirstArrival  int            `json:"first_arrival"`
	XFormula      string         `json:"x_formula"`
	DefragFile    string         `json:"defrag_events_file"`
	DefragDType   string         `json:"defrag_events_dtype"`
	DefragColumns []string       `json:"defrag_events_columns"`
	DefragEvents  int64          `json:"defrag_events"`
	StartedAt     string         `json:"started_at"`
	FinishedAt    string         `json:"finished_at"`
	Run           map[string]any `json:"run"`
}

func newFragSink(prefix string, stride int, seriesNames []string, extra map[string]any) (*fragSink, error) {
	if prefix == "" {
		return nil, fmt.Errorf("fragmentation output prefix is empty")
	}
	if stride < 1 {
		stride = 1
	}
	if dir := filepath.Dir(prefix); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("creating fragmentation output directory: %w", err)
		}
	}

	dataFile, err := os.Create(prefix + ".f32")
	if err != nil {
		return nil, fmt.Errorf("creating fragmentation data file: %w", err)
	}
	defragFile, err := os.Create(prefix + ".defrag.u64")
	if err != nil {
		dataFile.Close()
		return nil, fmt.Errorf("creating defragmentation events file: %w", err)
	}

	return &fragSink{
		prefix:      prefix,
		stride:      stride,
		dataFile:    dataFile,
		data:        bufio.NewWriterSize(dataFile, fragBufferSize),
		defragFile:  defragFile,
		defrag:      bufio.NewWriterSize(defragFile, 1<<16),
		row:         make([]byte, 4*len(seriesNames)),
		seriesNames: seriesNames,
		extra:       extra,
		startedAt:   time.Now(),
	}, nil
}

// writeRow appends one sample; values must have len(seriesNames) entries.
func (f *fragSink) writeRow(values []float64) error {
	for i, v := range values {
		binary.LittleEndian.PutUint32(f.row[i*4:], math.Float32bits(float32(v)))
	}
	if _, err := f.data.Write(f.row); err != nil {
		return err
	}
	f.samples++
	return nil
}

func (f *fragSink) writeDefrag(arrival, moved int) error {
	binary.LittleEndian.PutUint64(f.pair[0:], uint64(arrival))
	binary.LittleEndian.PutUint64(f.pair[8:], uint64(moved))
	if _, err := f.defrag.Write(f.pair[:]); err != nil {
		return err
	}
	f.defragCount++
	return nil
}

func (f *fragSink) close() error {
	errs := []error{f.data.Flush(), f.defrag.Flush(), f.dataFile.Close(), f.defragFile.Close()}
	for _, err := range errs {
		if err != nil {
			return fmt.Errorf("closing fragmentation output: %w", err)
		}
	}

	meta := fragMeta{
		FormatVersion: fragFormatVersion,
		DType:         "<f4",
		Layout:        "row-major [sample][series]",
		DataFile:      filepath.Base(f.prefix + ".f32"),
		NSeries:       len(f.seriesNames),
		SeriesNames:   f.seriesNames,
		NSamples:      f.samples,
		Stride:        f.stride,
		FirstArrival:  f.stride,
		XFormula:      "arrival = (sample_index + 1) * stride",
		DefragFile:    filepath.Base(f.prefix + ".defrag.u64"),
		DefragDType:   "<u8",
		DefragColumns: []string{"arrival", "connections_moved"},
		DefragEvents:  f.defragCount,
		StartedAt:     f.startedAt.Format(time.RFC3339),
		FinishedAt:    time.Now().Format(time.RFC3339),
		Run:           f.extra,
	}
	raw, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding fragmentation metadata: %w", err)
	}
	if err := os.WriteFile(f.prefix+".meta.json", raw, 0o644); err != nil {
		return fmt.Errorf("writing fragmentation metadata: %w", err)
	}
	return nil
}
