package helpers

import (
	"encoding/json"
	"log"
	"os"
	"simulator/connections"
	"simulator/infrastructure"
)

func ReadCapacityFile(capacityPath string) (infrastructure.Capacity, error) {
	dataBytes, err := os.ReadFile(capacityPath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return infrastructure.Capacity{}, err
	}

	var capacities infrastructure.Capacity

	err = json.Unmarshal(dataBytes, &capacities)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return infrastructure.Capacity{}, err
	}

	for i := range capacities.Bands {
		capacities.Bands[i].Slots = make([]bool, capacities.Bands[i].Capacity)
	}

	return capacities, nil
}

func ReadNetworkFile(networkPath string) (infrastructure.Network, error) {
	dataBytesNetwork, err := os.ReadFile(networkPath)

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return infrastructure.Network{}, err
	}

	var network infrastructure.Network

	err = json.Unmarshal(dataBytesNetwork, &network)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return infrastructure.Network{}, err
	}

	return network, nil
}

func ReadBitRateFile(bitRatePath string) (connections.BitRate, error) {
	dataBytesBitrate, err := os.ReadFile(bitRatePath)

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return connections.BitRate{}, err
	}

	var bitrate connections.BitRate

	err = json.Unmarshal(dataBytesBitrate, &bitrate)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return connections.BitRate{}, err
	}

	return bitrate, nil
}

func NetworkGenerate(networkPath string, capacityPath string, bitRatePath string) (infrastructure.Network, error) {

	capacities, err := ReadCapacityFile(capacityPath)

	if err != nil {
		log.Fatalf("Error reading capacities file: %v", err)
	}

	bitRate, err := ReadBitRateFile(bitRatePath)

	if err != nil {
		log.Fatalf("Error reading bit rate file: %v", err)
	}

	network, err := ReadNetworkFile(networkPath)

	if err != nil {
		log.Fatalf("Error reading network file: %v", err)
	}

	return network, nil
}
