package connections

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/Kayres21/optical-mb-sim-go/pkg/validator"
)

// BitRate represents a single traffic magnitude (Gbps) and every modulation
// format supported for it, mirroring flex-net-sim-python's BitRate class:
// one entry per magnitude, with parallel per-modulation slices.
type BitRate struct {
	Value        float64
	Modulation   []string
	Slots        []int
	Reach        []float64
	Bands        [][]string
	SlotsPerBand [][]int
	ReachPerBand [][]float64
}

// BitRateList is the flat list of BitRate magnitudes read from a bitrate
// file, indexed directly by the simulator's uniform bitrate selection.
type BitRateList struct {
	BitRates []BitRate
}

// AddModulation appends a modulation format to the BitRate, mirroring
// BitRate.add_modulation in the Python library. band/slotsPerBand/reachPerBand
// are only recorded when non-nil (multi-band files).
func (b *BitRate) AddModulation(modulation string, slots int, reach float64, band []string, slotsPerBand []int, reachPerBand []float64) {
	b.Modulation = append(b.Modulation, modulation)
	b.Slots = append(b.Slots, slots)
	b.Reach = append(b.Reach, reach)
	if band != nil {
		b.Bands = append(b.Bands, band)
		b.SlotsPerBand = append(b.SlotsPerBand, slotsPerBand)
		b.ReachPerBand = append(b.ReachPerBand, reachPerBand)
	}
}

// NumberOfModulations returns how many modulation formats this BitRate has.
func (b *BitRate) NumberOfModulations() int {
	return len(b.Modulation)
}

// GetDistanceAdaptive mirrors BitRate.get_distance_adaptive: it returns the
// index of the modulation with the smallest reach that still covers length,
// or -1 if no modulation reaches that far.
func (b *BitRate) GetDistanceAdaptive(length float64) int {
	best := -1
	minReach := math.Inf(1)

	for i, reach := range b.Reach {
		if reach >= length && reach < minReach {
			minReach = reach
			best = i
		}
	}

	return best
}

// SlotsForBand returns the slot count required by the given modulation for
// the given band name, falling back to the modulation's aggregate slot count
// when no per-band breakdown is present (single-band bitrate files).
func (b *BitRate) SlotsForBand(modulation int, bandName string) int {
	if modulation < len(b.Bands) {
		for i, name := range b.Bands[modulation] {
			if name == bandName {
				return b.SlotsPerBand[modulation][i]
			}
		}
	}
	if modulation < len(b.Slots) {
		return b.Slots[modulation]
	}
	return 0
}

// ─── Standard (schema-validated) bitrate file format ─────────────────────────

type schemaBitRateFile struct {
	BitRates []schemaBitRate `json:"bitrates"`
}

type schemaBitRate struct {
	Modulation string        `json:"modulation"`
	Slots      []schemaSlots `json:"slots"`
	Reachs     []schemaReach `json:"reachs"`
}

type schemaSlots struct {
	Gigabits string `json:"gigabits"`
	Slots    int    `json:"slots"`
}

type schemaReach struct {
	NumberOfBands int                  `json:"number_of_bands"`
	ReachsPerBand []schemaReachPerBand `json:"reachs_per_band"`
}

type schemaReachPerBand struct {
	Band  string  `json:"band"`
	Reach float64 `json:"reach"`
}

// ReadBitRateFile reads the schema-validated "bitrates" JSON format and
// converts it into the flat, per-magnitude BitRateList used at runtime.
// numberOfBands selects which of the file's per-band-count reach entries to
// use, matching the simulator's configured number of bands.
func ReadBitRateFile(bitRatePath string, numberOfBands int) (BitRateList, error) {
	schemaPath := filepath.Join(filepath.Dir(bitRatePath), "schema.json")
	dataBytesBitrate, err := validator.ValidateFile(bitRatePath, schemaPath)
	if err != nil {
		return BitRateList{}, fmt.Errorf("validating bitrate file: %w", err)
	}

	var raw schemaBitRateFile
	if err = json.Unmarshal(dataBytesBitrate, &raw); err != nil {
		return BitRateList{}, fmt.Errorf("parsing bitrate file %q: %w", bitRatePath, err)
	}

	return convertSchemaBitRateFile(raw, numberOfBands), nil
}

// convertSchemaBitRateFile groups the modulation-keyed schema format into one
// BitRate per magnitude (matching the Python data model).
func convertSchemaBitRateFile(raw schemaBitRateFile, numberOfBands int) BitRateList {
	magnitudes := make(map[string]float64)
	for _, mod := range raw.BitRates {
		for _, s := range mod.Slots {
			if v, err := strconv.ParseFloat(s.Gigabits, 64); err == nil {
				magnitudes[s.Gigabits] = v
			}
		}
	}

	keys := make([]string, 0, len(magnitudes))
	for k := range magnitudes {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return magnitudes[keys[i]] < magnitudes[keys[j]] })

	var result BitRateList
	for _, key := range keys {
		br := BitRate{Value: magnitudes[key]}

		for _, mod := range raw.BitRates {
			slots, ok := findSchemaSlots(mod.Slots, key)
			if !ok {
				continue
			}

			reachEntry := selectSchemaReach(mod.Reachs, numberOfBands)
			reach := 0.0
			var bandNames []string
			var slotsPerBand []int
			var reachPerBand []float64
			for i, rpb := range reachEntry.ReachsPerBand {
				if i == 0 {
					reach = rpb.Reach
				}
				bandNames = append(bandNames, rpb.Band)
				slotsPerBand = append(slotsPerBand, slots)
				reachPerBand = append(reachPerBand, rpb.Reach)
			}

			br.AddModulation(mod.Modulation, slots, reach, bandNames, slotsPerBand, reachPerBand)
		}

		result.BitRates = append(result.BitRates, br)
	}

	return result
}

func findSchemaSlots(slots []schemaSlots, gigabits string) (int, bool) {
	for _, s := range slots {
		if s.Gigabits == gigabits {
			return s.Slots, true
		}
	}
	return 0, false
}

// selectSchemaReach picks the reach entry matching numberOfBands, falling
// back to the entry with the largest number_of_bands not exceeding it, or
// the first entry if none qualify.
func selectSchemaReach(reachs []schemaReach, numberOfBands int) schemaReach {
	best := schemaReach{}
	bestBands := -1
	for _, r := range reachs {
		if r.NumberOfBands == numberOfBands {
			return r
		}
		if r.NumberOfBands <= numberOfBands && r.NumberOfBands > bestBands {
			best = r
			bestBands = r.NumberOfBands
		}
	}
	if bestBands == -1 && len(reachs) > 0 {
		return reachs[0]
	}
	return best
}

// SelectBitrateMethod mirrors BitRate.select_bit_rate_method: an empty
// fileName yields the built-in default bitrates, otherwise the schema-based
// file is read and converted.
func SelectBitrateMethod(fileName string, numberOfBands int) (BitRateList, error) {
	if fileName == "" {
		return defaultBitRates(), nil
	}
	return ReadBitRateFile(fileName, numberOfBands)
}

func defaultBitRates() BitRateList {
	// Defaults: 10, 40, 100, 400, 1000 Gbps, all BPSK, matching the Python
	// library's default_bit_rates().
	defaults := []struct {
		value float64
		slots int
	}{
		{10, 1},
		{40, 4},
		{100, 8},
		{400, 32},
		{1000, 80},
	}

	var rates BitRateList
	for _, d := range defaults {
		br := BitRate{Value: d.value}
		br.AddModulation("BPSK", d.slots, 5520, nil, nil, nil)
		rates.BitRates = append(rates.BitRates, br)
	}
	return rates
}

func TrasnformIntToModulation(modulation int) string {
	switch modulation {
	case 0:
		return "BPSK"
	case 1:
		return "QPSK"
	case 2:
		return "8-QAM"
	case 3:
		return "16-QAM"
	default:
		log.Fatalf("Invalid modulation type: %d", modulation)
		return "BPSK" // Default case, should not be reached
	}
}
