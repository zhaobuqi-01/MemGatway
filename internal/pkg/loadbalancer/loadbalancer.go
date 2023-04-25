package loadbalancer

type LoadBalancer interface {
	// Add one or more backend servers to the load balancer.
	Add(servers ...string) error

	// Remove a specified backend server from the load balancer.
	Remove(server string) error

	// Get the selected backend server based on the load balancing algorithm.
	Get(clientID string) (string, error)

	// Update the backend servers, such as adding or removing servers,
	// or updating the weight or health status of a server.
	Update()
}
