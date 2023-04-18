// loadbalancer/consistenthashing_test.go

package loadbalancer

import (
	"testing"
)

// TestConsistentHashing_Add_Get tests the Add and Get methods of the ConsistentHashing struct.
func TestConsistentHashing_Add_Get(t *testing.T) {
	consistentHashing := NewConsistentHashing(100)
	servers := []string{
		"server1",
		"server2",
		"server3",
	}

	err := consistentHashing.Add(servers...)
	if err != nil {
		t.Fatalf("Failed to add servers: %v", err)
	}

	key := "test-key"
	server, err := consistentHashing.Get(key)
	if err != nil {
		t.Fatalf("Failed to get server: %v", err)
	}

	t.Logf("Server for key '%s': %s", key, server)
}

// TestConsistentHashing_Remove tests the Remove method of the ConsistentHashing struct.
func TestConsistentHashing_Remove(t *testing.T) {
	consistentHashing := NewConsistentHashing(100)
	servers := []string{
		"server1",
		"server2",
		"server3",
	}

	err := consistentHashing.Add(servers...)
	if err != nil {
		t.Fatalf("Failed to add servers: %v", err)
	}

	err = consistentHashing.Remove("server2")
	if err != nil {
		t.Fatalf("Failed to remove server: %v", err)
	}

	key := "test-key"
	server, err := consistentHashing.Get(key)
	if err != nil {
		t.Fatalf("Failed to get server: %v", err)
	}

	if server == "server2" {
		t.Errorf("Removed server 'server2' is still being returned by Get method")
	}

	t.Logf("Server for key '%s': %s", key, server)
}
