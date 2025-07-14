package connections

import (
	"encoding/json"
	"log"
	"os"
)

func ReadBitRateFile(bitRatePath string) (BitRate, error) {
	dataBytesBitrate, err := os.ReadFile(bitRatePath)

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return BitRate{}, err
	}

	var bitrate BitRate

	err = json.Unmarshal(dataBytesBitrate, &bitrate)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return BitRate{}, err
	}

	return bitrate, nil
}
