package infrastructure

import (
	"encoding/json"
	"log"
	"os"
)

type Network struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

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

func cloneCapacity(orig Capacity) Capacity {
	bands := make([]Band, len(orig.Bands))
	for i, b := range orig.Bands {
		slots := make([]bool, len(b.Slots))
		copy(slots, b.Slots)
		bands[i] = Band{Slots: slots}
	}
	return Capacity{Bands: bands}
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
		network.Links[i].Capacities = cloneCapacity(capacities)
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

func (n *Network) GetLinkByPath(ids []int) []*Link {
	links := make([]*Link, 0, len(ids))
	for i, id := range ids {

		if i+1 < len(ids) {
			src := id
			dst := ids[i+1]
			link := n.GetLinkBySourceDestination(src, dst)
			if link != nil {
				links = append(links, link)
			}
		}

	}
	return links
}

func (n *Network) GetLinkBySourceDestination(src, dst int) *Link {
	for i := range n.Links {
		if n.Links[i].Source == src && n.Links[i].Destination == dst {
			return &n.Links[i]
		}
	}
	return nil
}
