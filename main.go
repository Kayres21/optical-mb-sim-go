package main

import (
	"fmt"
	"simulator/internal/connections"
	randomvariable "simulator/internal/connections/random_variable"
	"simulator/internal/infrastructure"
)

func main() {

	network, err := infrastructure.NetworkGenerate("files/networks/UKNet_BDM.json", "files/capacities/capacities.json")
	if err != nil {
		fmt.Printf("Error reading network file: %v\n", err)
		return
	}
	fmt.Println("Network Name:", network.Name)

	for _, link := range network.Links {
		fmt.Printf("Link ID: %d, Source: %d, Destination: %d, Length: %d\n", link.ID, link.Source, link.Destination, link.Length)
		for _, band := range link.Capacities.Bands {
			fmt.Printf("  Band ID: %s, Name: %s, Capacity: %d\n", band.ID, band.Name, band.Capacity)
		}
	}

	bitRate, err := connections.ReadBitRateFile("files/bitrate/bitrate.json")

	if err != nil {
		fmt.Printf("Error reading bitrate file: %v\n", err)
		return
	}

	for _, br := range bitRate.BitRates {
		fmt.Printf("Bitrate Modulation: %s\n", br.Modulation)
		for _, slot := range br.Slots {
			fmt.Printf("  Gigabits: %s, Slots: %d\n", slot.Gigabits, slot.Slots)
		}
		for _, reach := range br.Reachs {
			fmt.Printf("  Number of Bands: %d\n", reach.NumberOfBands)
			for _, rpb := range reach.Reach {
				fmt.Printf("    Band: %s, Reach: %d\n", rpb.Band, rpb.Reach)
			}
		}
	}

	var RandomVariable randomvariable.RandomVariable

	RandomVariable.Arrive.SetParameter(10)
	RandomVariable.Departure.SetParameter(2)
	RandomVariable.BitrateSelect.SetParameter(3)
	RandomVariable.SourceNodeSelect.SetParameter(30)
	RandomVariable.DestinationNodeSelect.SetParameter(30)
	RandomVariable.BandSelect.SetParameter(4)

	RandomVariable.SetSeedArrive(1)
	RandomVariable.SetSeedDeparture(2)
	RandomVariable.SetSeedBitrate(3)
	RandomVariable.SetSeedSource(4)
	RandomVariable.SetSeedDestination(5)
	RandomVariable.SetSeedBand(6)

	fmt.Printf("Arrival Random Variable: %f\n", RandomVariable.GetNetValueExponential("arrive"))
	fmt.Printf("Arrival Random Variable: %f\n", RandomVariable.GetNetValueExponential("arrive"))
	fmt.Printf("Departure Random Variable: %f\n", RandomVariable.GetNetValueExponential("departure"))
	fmt.Printf("Bitrate Random Variable: %d\n", RandomVariable.GetNetValueUniform("bitrate"))
	fmt.Printf("Source Node Random Variable: %d\n", RandomVariable.GetNetValueUniform("source"))
	fmt.Printf("Destination Node Random Variable: %d\n", RandomVariable.GetNetValueUniform("destination"))
	fmt.Printf("Band Random Variable: %d\n", RandomVariable.GetNetValueUniform("band"))
	fmt.Printf("Band Random Variable: %d\n", RandomVariable.GetNetValueUniform("band"))

}
