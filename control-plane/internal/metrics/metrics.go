package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cp_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "code"},
	)
)

func init() {
	prometheus.MustRegister(RequestCounter)
}