package connections

import (
	"encoding/json"
	"log"
	"os"
)

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
