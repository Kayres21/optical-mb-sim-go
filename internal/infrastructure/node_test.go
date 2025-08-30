package infrastructure

import "testing"

func TestGetID(t *testing.T) {
	node := Node{ID: 1}
	if node.GetID() != 1 {
		t.Errorf("Expected ID to be 1, got %d", node.GetID())
	}
}
