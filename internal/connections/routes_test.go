package connections

import "testing"

func TestReadRoutesFile(t *testing.T) {
	routesFile := "../../files/routes/network_test_routes.json"
	routes, err := ReadRoutesFile(routesFile)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if routes.Name != "network_test_routes" {
		t.Errorf("Expected routes name to be 'network_test_routes', got '%s'", routes.Name)
	}

	if routes.Alias != "Test Routes" {
		t.Errorf("Expected routes alias to be 'Test Routes', got '%s'", routes.Alias)
	}

	if len(routes.Paths) != 6 {
		t.Errorf("Expected 6 paths, got %d", len(routes.Paths))
	}
}

func TestGetKshortestPath(t *testing.T) {
	routesFile := "../../files/routes/network_test_routes.json"
	routes, err := ReadRoutesFile(routesFile)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	path := routes.GetKshortestPath(0, 0, 2)
	expectedPath := []int{0, 2}
	if len(path) != len(expectedPath) {
		t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
	} else {
		for i := range path {
			if path[i] != expectedPath[i] {
				t.Errorf("Expected path %v, got %v", expectedPath, path)
				break
			}
		}
	}
}
