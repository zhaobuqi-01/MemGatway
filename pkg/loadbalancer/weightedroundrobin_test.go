// loadbalancer/weightedroundrobin_test.go

package loadbalancer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWeightRoundRobinBalance_LoadBalancing(t *testing.T) {
	wrr := &WeightRoundRobinBalance{}

	// Add three nodes with different weights
	err := wrr.Add("127.0.0.1:8001", "5")
	assert.NoError(t, err)
	err = wrr.Add("127.0.0.1:8002", "2")
	assert.NoError(t, err)
	err = wrr.Add("127.0.0.1:8003", "3")
	assert.NoError(t, err)

	// Set up a map to store the number of times each node is selected
	counts := make(map[string]int)

	// Select nodes 1000 times and count how many times each node is selected
	for i := 0; i < 1000; i++ {
		addr := wrr.get()
		counts[addr]++
	}

	// Calculate the percentage of times each node is selected
	pct1 := float64(counts["127.0.0.1:8001"]) / 1000
	pct2 := float64(counts["127.0.0.1:8002"]) / 1000
	pct3 := float64(counts["127.0.0.1:8003"]) / 1000

	// Check that the load balancing effect is within a reasonable range
	assert.InDelta(t, pct1, 0.5, 0.1)
	assert.InDelta(t, pct2, 0.2, 0.1)
	assert.InDelta(t, pct3, 0.3, 0.1)
}

// Test the Remove() method by adding three nodes, removing one node, and checking that the correct node is removed
func TestWeightRoundRobinBalance_Remove(t *testing.T) {
	wrr := &WeightRoundRobinBalance{}

	// Add three nodes with different weights
	err := wrr.Add("127.0.0.1:8001", "5")
	assert.NoError(t, err)
	err = wrr.Add("127.0.0.1:8002", "2")
	assert.NoError(t, err)
	err = wrr.Add("127.0.0.1:8003", "3")
	assert.NoError(t, err)

	// Remove one node and check that the correct node is removed
	err = wrr.Remove("127.0.0.1:8002")
	assert.NoError(t, err)
	assert.Equal(t, len(wrr.rss), 2)
	assert.Equal(t, wrr.rss[0].addr, "127.0.0.1:8001")
	assert.Equal(t, wrr.rss[1].addr, "127.0.0.1:8003")
}
