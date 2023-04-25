// loadbalancer/roundrobin.go

package loadbalancer

import (
	"errors"
	"sync"
)

// RoundRobin is a load balancing algorithm that selects servers in a circular order.
type RoundRobin struct {
	mu      sync.Mutex
	servers []string
	index   int
}

// NewRoundRobin initializes a new RoundRobin instance.
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		servers: make([]string, 0),
		index:   -1,
	}
}

// Add appends the provided servers to the RoundRobin instance's servers list.
func (r *RoundRobin) Add(servers ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.servers = append(r.servers, servers...)
	return nil
}

// Get returns the next server from the RoundRobin instance's servers list in a circular order.
func (r *RoundRobin) Get(_ string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.servers) == 0 {
		return "", errors.New("no available servers")
	}

	r.index = (r.index + 1) % len(r.servers)
	return r.servers[r.index], nil
}

// Remove removes the specified server from the RoundRobin instance's servers list.
func (r *RoundRobin) Remove(server string) error {
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

// Update can be used to update the servers list based on service discovery.
func (r *RoundRobin) Update() {
	// This method can be implemented to update the servers list when needed.
}
