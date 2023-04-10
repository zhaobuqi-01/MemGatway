// loadbalancer/consistenthashing.go

package loadbalancer

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// ConsistentHashing is a load balancing algorithm that uses consistent hashing for distributing keys among servers.
type ConsistentHashing struct {
	mu       sync.Mutex
	servers  map[uint32]string
	keys     []uint32
	replicas int
}

// NewConsistentHashing initializes a new ConsistentHashing instance with the given number of replicas.
func NewConsistentHashing(replicas int) *ConsistentHashing {
	return &ConsistentHashing{
		servers:  make(map[uint32]string),
		keys:     make([]uint32, 0),
		replicas: replicas,
	}
}

// Add appends the provided servers to the ConsistentHashing instance's servers list.
func (c *ConsistentHashing) Add(servers ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, server := range servers {
		for i := 0; i < c.replicas; i++ {
			hash := c.hash(server, i)
			c.servers[hash] = server
			c.keys = append(c.keys, hash)
		}
	}

	sort.Slice(c.keys, func(i, j int) bool {
		return c.keys[i] < c.keys[j]
	})

	return nil
}

// Get returns the server that is responsible for the given key.
func (c *ConsistentHashing) Get(key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.servers) == 0 {
		return "", errors.New("no available servers")
	}

	hash := crc32.ChecksumIEEE([]byte(key))
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hash
	})

	if idx == len(c.keys) {
		idx = 0
	}

	return c.servers[c.keys[idx]], nil
}

// Remove removes the specified server from the ConsistentHashing instance's servers list.
func (c *ConsistentHashing) Remove(server string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := 0; i < c.replicas; i++ {
		hash := c.hash(server, i)
		delete(c.servers, hash)

		for idx, key := range c.keys {
			if key == hash {
				c.keys = append(c.keys[:idx], c.keys[idx+1:]...)
				break
			}
		}
	}

	return nil
}

// Update can be used to update the servers list based on service discovery.
func (c *ConsistentHashing) Update() {
	// This method can be implemented to update the servers list when needed.
}

// hash generates a hash value for the given server and replica index.
func (c *ConsistentHashing) hash(server string, index int) uint32 {
	data := []byte(server + "-" + strconv.Itoa(index))
	return crc32.ChecksumIEEE(data)
}
