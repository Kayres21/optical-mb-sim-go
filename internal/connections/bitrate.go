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
	Gigabits string `json:"gigabits"`
	Slots    int    `json:"slots"`
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
