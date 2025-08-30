package infrastructure

import "testing"

func TestReadNetworkFile(t *testing.T) {
	networkPath := "../../files/networks/network_test.json"
	network, err := ReadNetworkFile(networkPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if network.Name != "Test Network" {
		t.Errorf("Expected network name to be 'Test Network', got '%s'", network.Name)
	}
	if len(network.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(network.Nodes))
	}
	if len(network.Links) != 1 {
		t.Errorf("Expected 1 link, got %d", len(network.Links))
	}
}
