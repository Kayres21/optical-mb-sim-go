package main

import (
	"fmt"
	"simulator/internal/infrastructure"
)

func main() {

	// Initialize the network
	network, err := infrastructure.NetworkGenerate("files/networks/UKNet_BDM.json", "files/capacities/capacities.json")
	if err != nil {
		panic(err)
	}
	fmt.Println("Network Name:", network.Name)

	for _, link := range network.Links {
		fmt.Printf("Link ID: %d, Source: %d, Destination: %d, Length: %d\n", link.ID, link.Source, link.Destination, link.Length)
		for _, band := range link.Capacities.Bands {
			fmt.Printf("  Band ID: %s, Name: %s, Capacity: %d\n", band.ID, band.Name, band.Capacity)
		}
	}
}
