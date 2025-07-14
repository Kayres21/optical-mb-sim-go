package infrastructure

import (
	"encoding/json"
	"log"
	"os"
)

func ReadNetworkFile(networkPath string) (Network, error) {
	dataBytesNetwork, err := os.ReadFile(networkPath)

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return Network{}, err
	}

	var network Network

	err = json.Unmarshal(dataBytesNetwork, &network)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return Network{}, err
	}

	return network, nil
}

func NetworkGenerate(networkPath string, capacityPath string) (Network, error) {

	capacities, err := ReadCapacityFile(capacityPath)

	if err != nil {
		log.Fatalf("Error reading capacities file: %v", err)
	}

	network, err := ReadNetworkFile(networkPath)

	if err != nil {
		log.Fatalf("Error reading network file: %v", err)
	}

	for i := range network.Links {
		network.Links[i].Capacities = capacities
	}

	return network, nil
}

func (n *Network) GetNodeByID(id int) *Node {
	for i := range n.Nodes {
		if n.Nodes[i].ID == id {
			return &n.Nodes[i]
		}
	}
	return nil
}
func (n *Network) GetLinkByID(id int) *Link {
	for i := range n.Links {
		if n.Links[i].ID == id {
			return &n.Links[i]
		}
	}
	return nil
}
