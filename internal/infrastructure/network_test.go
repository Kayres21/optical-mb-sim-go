package infrastructure

import "testing"

func TestReadNetworkFile(t *testing.T) {
	networkPath := "../../files/networks/network_test.json"
	network, err := ReadNetworkFile(networkPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if network.GetName() != "network_test" {
		t.Errorf("Expected network name to be 'Test Network', got '%s'", network.GetName())
	}
	if len(network.GetNodes()) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(network.GetNodes()))
	}
	if len(network.GetLinks()) != 6 {
		t.Errorf("Expected 6 link, got %d", len(network.GetLinks()))
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

		if !(link.GetSource() == 1 && link.GetDestination() == 2) {
			t.Errorf("Expected link between nodes 1 and 2, got link between %d and %d", link.GetSource(), link.GetDestination())
		}
	}
}
