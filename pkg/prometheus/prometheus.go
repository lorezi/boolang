package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	registry *prometheus.Registry
	handler  http.HandlerFunc
}

func New(custom bool) Prometheus {
	reg := prometheus.NewRegistry()

	if custom {
		return Prometheus{
			registry: reg,
			handler: func(w http.ResponseWriter, r *http.Request) {
				promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
			},
		}
	} else {
		return Prometheus{
			registry: reg,
			handler: func(w http.ResponseWriter, r *http.Request) {
				promhttp.Handler().ServeHTTP(w, r)
			},
		}
	}
}

// Handler returns HTTP handler
func (p Prometheus) Handler() http.HandlerFunc {
	return p.handler
}

// Registry returns `Registry` instance which helps registering the custom metric collections
func (p Prometheus) Registry() *prometheus.Registry {
	return p.registry
}
