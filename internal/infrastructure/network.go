package infrastructure

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/Kayres21/optical-mb-sim-go/pkg/validator"
)

type Network struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func ReadNetworkFile(networkPath string) (Network, error) {
	schemaPath := filepath.Join(filepath.Dir(networkPath), "schema.json")
	data, err := validator.ValidateFile(networkPath, schemaPath)
	if err != nil {
		return Network{}, fmt.Errorf("validating network file: %w", err)
	}

	var network Network
	if err = json.Unmarshal(data, &network); err != nil {
		return Network{}, fmt.Errorf("parsing network file %q: %w", networkPath, err)
	}

	return network, nil
}

func cloneCapacity(orig Capacity) Capacity {
	bands := make([]Band, len(orig.Bands))
	for i, b := range orig.Bands {
		slots := make([]bool, len(b.Slots))
		copy(slots, b.Slots)
		bands[i] = Band{
			ID:       b.ID,
			Name:     b.Name,
			SlotsLen: b.SlotsLen,
			Slots:    slots,
		}
	}
	return Capacity{Bands: bands}
}

func NetworkGenerate(networkPath string, capacityPath string) (Network, error) {
	capacities, err := ReadCapacityFile(capacityPath)
	if err != nil {
		return Network{}, fmt.Errorf("generating network: %w", err)
	}

	network, err := ReadNetworkFile(networkPath)
	if err != nil {
		return Network{}, fmt.Errorf("generating network: %w", err)
	}

	for i := range network.Links {
		network.Links[i].Capacities = cloneCapacity(capacities)
		network.Links[i].UpdateAllFragmentationRatios()
	}

	return network, nil
}

func (n Network) Clone() Network {
	clone := n
	clone.Nodes = append([]Node(nil), n.Nodes...)
	clone.Links = make([]Link, len(n.Links))
	for i := range n.Links {
		clone.Links[i] = n.Links[i]
		clone.Links[i].Capacities = cloneCapacity(n.Links[i].Capacities)
		clone.Links[i].FragmentationRatioByBand = append([]float64(nil), n.Links[i].FragmentationRatioByBand...)
	}
	return clone
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
			link := n.GetLinkBySourceDestination(id, ids[i+1])
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

func (n *Network) GetPathDistance(links []*Link) float64 {
	distance := 0.0
	for _, link := range links {
		distance += link.Length
	}
	return distance
}

// FragmentationRatio returns the mean of the selected bands' average FR values.
func (n *Network) FragmentationRatio(numberOfBands int) float64 {
	ratios := n.FragmentationRatiosByBand(numberOfBands)
	if len(ratios) == 0 {
		return 0
	}

	total := 0.0
	for _, ratio := range ratios {
		total += ratio
	}
	return total / float64(len(ratios))
}

// FragmentationRatiosByBand returns the average FR for every selected band
// across all links in the network.
func (n *Network) FragmentationRatiosByBand(numberOfBands int) []float64 {
	if len(n.Links) == 0 || numberOfBands <= 0 {
		return nil
	}

	bandCount := numberOfBands
	availableBands := len(n.Links[0].Capacities.Bands)
	if availableBands == 0 {
		availableBands = len(n.Links[0].FragmentationRatioByBand)
	}
	if bandCount > availableBands {
		bandCount = availableBands
	}

	ratios := make([]float64, bandCount)
	for i := range n.Links {
		for band := range ratios {
			ratios[band] += n.Links[i].GetFragmentationRatioByBand(band)
		}
	}
	for band := range ratios {
		ratios[band] /= float64(len(n.Links))
	}

	return ratios
}
