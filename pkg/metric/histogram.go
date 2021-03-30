package metric

import "github.com/prometheus/client_golang/prometheus"

func ResponseDurationHistogram() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "client",
		Name:      "balance_response_duration_histogram",
		Help:      "Balance response duration (ms)",
		Buckets:   []float64{10, 50, 90, 130, 170, 210, 250, 290, 330},
	}, []string{"operation"})
}
