// loadbalancer/random_test.go

package loadbalancer

import (
	"testing"
)

// TestRandom_Get tests the Get method of the Random struct.
func TestRandom_Get(t *testing.T) {
	random := NewRandom()
	servers := []string{"server1", "server2", "server3"}

	// Add servers to the Random instance
	err := random.Add(servers...)
	if err != nil {
		t.Fatalf("Failed to add servers: %v", err)
	}

	// Test the random algorithm
	for i := 0; i < 10; i++ {
		server, err := random.Get("")
		if err != nil {
			t.Fatalf("Failed to get server: %v", err)
		}

		// Ensure that the selected server is in the servers list
		found := false
		for _, s := range servers {
			if server == s {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Selected server %s is not in the servers list", server)
		}
	}
}

// TestRandom_Remove tests the Remove method of the Random struct.
func TestRandom_Remove(t *testing.T) {
	random := NewRandom()
	servers := []string{"server1", "server2", "server3"}

	// Add servers to the Random instance
	err := random.Add(servers...)
	if err != nil {
		t.Fatalf("Failed to add servers: %v", err)
	}

	// Remove a server from the Random instance
	err = random.Remove("server2")
	if err != nil {
		t.Fatalf("Failed to remove server: %v", err)
	}

	// Ensure that the server was removed
	for _, server := range random.servers {
		if server == "server2" {
			t.Errorf("Server 'server2' was not removed")
		}
	}
}
