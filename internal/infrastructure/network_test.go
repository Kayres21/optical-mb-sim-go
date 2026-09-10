package infrastructure

import "testing"

func TestReadNetworkFile(t *testing.T) {
	networkPath := "../../files/networks/network_test.json"
	network, err := ReadNetworkFile(networkPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if network.Name != "network_test" {
		t.Errorf("Expected network name to be 'Test Network', got '%s'", network.Name)
	}
	if len(network.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(network.Nodes))
	}
	if len(network.Links) != 6 {
		t.Errorf("Expected 6 link, got %d", len(network.Links))
	}
}

func TestGetLinkByPath(t *testing.T) {
	networkPath := "../../files/networks/network_test.json"
	network, err := ReadNetworkFile(networkPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	ids := [2]int{1, 2}
	links := network.GetLinkByPath(ids[:])
	if links == nil {
		t.Fatalf("Expected to find link between nodes 1 and 2, got nil")
	}

	for _, link := range links {

		if !(link.Source == 1 && link.Destination == 2) {
			t.Errorf("Expected link between nodes 1 and 2, got link between %d and %d", link.Source, link.Destination)
		}
	}
}

func TestNetworkFragmentationRatio(t *testing.T) {
	network := Network{
		Links: []Link{
			{FragmentationRatioByBand: []float64{0.2, 0.4}},
			{FragmentationRatioByBand: []float64{0.6, 0.8}},
		},
	}

	if got := network.FragmentationRatio(1); got != 0.4 {
		t.Fatalf("expected one-band fragmentation ratio 0.4, got %v", got)
	}
	if got := network.FragmentationRatio(2); got != 1.0 {
		t.Fatalf("expected two-band fragmentation ratio 1.0, got %v", got)
	}
}
