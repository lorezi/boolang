package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

func HttpResponseCounter() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "client",
		Name:      "http_response_counter",
		Help:      "Number of HTTP responses",
	}, []string{"operation", "code"})
}

func BalanceActivityCounter() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "client",
		Name:        "balance_activity_counter",
		Help:        "Balance activity history",
		ConstLabels: map[string]string{},
	}, []string{"activity", "client"})
}
