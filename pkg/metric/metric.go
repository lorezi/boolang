package metric

import "github.com/prometheus/client_golang/prometheus"

type Metric struct {
	HttpResponseCounter       *prometheus.CounterVec
	BalanceActivityCounter    *prometheus.CounterVec
	BalanceGauge              prometheus.Gauge
	ResponseDurationHistogram *prometheus.HistogramVec
}

func New(registry *prometheus.Registry) Metric {
	m := &Metric{}

	m.HttpResponseCounter = HttpResponseCounter()
	registry.MustRegister(m.HttpResponseCounter)

	m.BalanceActivityCounter = BalanceActivityCounter()
	registry.MustRegister(m.BalanceActivityCounter)

	m.BalanceGauge = BalanceGauge()
	registry.MustRegister(m.BalanceGauge)

	m.ResponseDurationHistogram = ResponseDurationHistogram()
	registry.MustRegister(m.ResponseDurationHistogram)

	return *m

}
