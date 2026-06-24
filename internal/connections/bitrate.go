package connections

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/Kayres21/optical-mb-sim-go/pkg/validator"
)

type BitRateList struct {
	BitRates []BitRate `json:"bitrates"`
}

type BitRate struct {
	Modulation string  `json:"modulation"`
	Slots      []Slots `json:"slots"`
	Reachs     []Reach `json:"reachs"`
}

type Slots struct {
	Gigabits     string         `json:"gigabits"`
	Slots        int            `json:"slots"`
	SlotsPerBand map[string]int `json:"slots_per_band,omitempty"`
}

type Reach struct {
	NumberOfBands int            `json:"number_of_bands"`
	Reach         []ReachPerBand `json:"reachs_per_band"`
}

type ReachPerBand struct {
	Band  string `json:"band"`
	Reach int    `json:"reach"`
}

func ReadBitRateFile(bitRatePath string) (BitRateList, error) {
	schemaPath := filepath.Join(filepath.Dir(bitRatePath), "schema.json")
	dataBytesBitrate, err := validator.ValidateFile(bitRatePath, schemaPath)
	if err != nil {
		return BitRateList{}, fmt.Errorf("validating bitrate file: %w", err)
	}

	var bitrate BitRateList

	if err = json.Unmarshal(dataBytesBitrate, &bitrate); err != nil {
		return BitRateList{}, fmt.Errorf("parsing bitrate file %q: %w", bitRatePath, err)
	}

	return bitrate, nil
}

// SelectBitrateMethod mirrors the C++ BitRate::selectBitrateMethod behaviour:
// - If fileName is empty, return a set of default bitrates.
// - Otherwise read the bitrate JSON file.
func SelectBitrateMethod(fileName string, networkType int) (BitRateList, error) {
	if fileName == "" {
		return defaultBitRates(), nil
	}
	// For now, networkType is not used to change parsing behaviour; rely on
	// ReadBitRateFile which supports both structures via JSON schema validation.
	return ReadBitRateFile(fileName)
}

func defaultBitRates() BitRateList {
	// Defaults: 10, 40, 100, 400, 1000 with slots 1,4,8,32,80 and BPSK modulation
	rates := BitRateList{}
	defaults := []struct {
		gb   string
		slots int
	}{
		{"10", 1},
		{"40", 4},
		{"100", 8},
		{"400", 32},
		{"1000", 80},
	}

	for _, d := range defaults {
		br := BitRate{Modulation: "BPSK"}
		br.Slots = []Slots{{Gigabits: d.gb, Slots: d.slots}}
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
