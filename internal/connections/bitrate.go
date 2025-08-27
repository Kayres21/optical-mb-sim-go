package connections

import (
	"encoding/json"
	"log"
	"os"
)

type BitRateList struct {
	BitRates []BitRate `json:"bitrates"`
}

type BitRate struct {
	Modulation string  `json:"modulation"`
	Slots      []Slots `json:"slots"`
	Reachs     []Reach `json:"reach"`
}

type Slots struct {
	Gigabits string `json:"gigabits"`
	Slots    int    `json:"slots"`
}

type Reach struct {
	NumberOfBands int            `json:"number_of_bands"`
	Reach         []ReachPerBand `json:"reach"`
}

type ReachPerBand struct {
	Band  string `json:"band"`
	Reach int    `json:"reach"`
}

func ReadBitRateFile(bitRatePath string) (BitRateList, error) {
	dataBytesBitrate, err := os.ReadFile(bitRatePath)

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return BitRateList{}, err
	}

	var bitrate BitRateList

	err = json.Unmarshal(dataBytesBitrate, &bitrate)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return BitRateList{}, err
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
