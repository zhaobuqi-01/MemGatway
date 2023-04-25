package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	AppQPS = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_qps",
			Help: "QPS per App",
		},
		[]string{"app_id"},
	)
	AppQPD = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_qpd",
			Help: "QPD per App",
		},
		[]string{"app_id"},
	)
	LoadBalancerCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "load_balancer_total",
			Help: "Total number of requests processed by the load balancer.",
		},
		[]string{"service_name", "load_balancer_type"},
	)
	TotalTrafficCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_traffic_bytes",
			Help: "Total traffic in bytes processed by the gateway.",
		},
		[]string{"direction"}, // "inbound" or "outbound"
	)
)

func init() {
	prometheus.MustRegister(LoadBalancerCounter)
	prometheus.MustRegister(TotalTrafficCounter)
	prometheus.MustRegister(AppQPS)
	prometheus.MustRegister(AppQPD)
}
