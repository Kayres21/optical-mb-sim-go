package helpers

import (
	"encoding/json"
	"log"
	"os"
	"simulator/connections"
	"simulator/network"
)

func ReadCapacityFile(capacityPath string) network.Capacity {
	dataBytes, err := os.ReadFile(capacityPath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	var capacities network.Capacity

	err = json.Unmarshal(dataBytes, &capacities)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	for i := range capacities.Bands {
		capacities.Bands[i].Slots = make([]bool, capacities.Bands[i].Capacity)
	}

	return capacities
}

func ReadNetworkFile(networkPath string) network.Network {
	dataBytesNetwork, err := os.ReadFile(networkPath)

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	var network network.Network

	err = json.Unmarshal(dataBytesNetwork, &network)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	return network
}

func ReadBitRateFile(bitRatePath string) connections.BitRate {
	dataBytesBitrate, err := os.ReadFile(bitRatePath)

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	var bitrate connections.BitRate

	err = json.Unmarshal(dataBytesBitrate, &bitrate)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	return bitrate
}

func NetworkGenerate(networkPath string, capacityPath string, bitRatePath string) network.Network {

	capacities := ReadCapacityFile(capacityPath)

	bitRate := ReadBitRateFile(bitRatePath)

	network := ReadNetworkFile(networkPath)

	return network
}
