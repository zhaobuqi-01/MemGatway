// loadbalancer/random.go

package loadbalancer

import (
	"errors"
	"math/rand"
	"sync"
)

// Random is a load balancing algorithm that selects servers randomly.
type Random struct {
	mu      sync.Mutex
	servers []string
}

// NewRandom initializes a new Random instance.
func NewRandom() *Random {
	return &Random{
		servers: make([]string, 0),
	}
}

// Add appends the provided servers to the Random instance's servers list.
func (r *Random) Add(servers ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.servers = append(r.servers, servers...)
	return nil
}

// Remove removes the specified server from the Random instance's servers list.
func (r *Random) Remove(server string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, s := range r.servers {
		if s == server {
			r.servers = append(r.servers[:i], r.servers[i+1:]...)
			return nil
		}
	}
	return errors.New("server not found")
}

// Get returns a random server from the Random instance's servers list.
func (r *Random) Get(_ string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.servers) == 0 {
		return "", errors.New("no available servers")
	}

	index := rand.Intn(len(r.servers))
	return r.servers[index], nil
}

// Update can be used to update the servers list based on service discovery.
func (r *Random) Update() {
	// This method can be implemented to update the servers list when needed.
}
