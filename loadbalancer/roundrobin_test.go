// loadbalancer/roundrobin_test.go

package loadbalancer

import (
	"testing"
)

// TestRoundRobin_Get tests the Get method of the RoundRobin struct.
func TestRoundRobin_Get(t *testing.T) {
	rr := NewRoundRobin()
	nodes := []string{"node1", "node2", "node3"}

	// Add nodes to the RoundRobin instance
	err := rr.Add(nodes...)
	if err != nil {
		t.Fatalf("Failed to add nodes: %v", err)
	}

	// Test the round-robin algorithm
	for i := 0; i < 6; i++ {
		node, err := rr.Get("")
		if err != nil {
			t.Fatalf("Failed to get node: %v", err)
		}

		expectedNode := nodes[i%len(nodes)]
		if node != expectedNode {
			t.Errorf("Expected node %s, but got %s", expectedNode, node)
		}
	}
}
