package farawaymetrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "faraway"

// MetricsServer holds configuration for our metrics server
type MetricsServer struct {
	Registry      *prometheus.Registry // Adjusted to *prometheus.Registry
	ServiceHealth *prometheus.GaugeVec
}

// NewMetricsServer initializes a new MetricsServer with default collectors and custom metrics
func NewMetricsServer() *MetricsServer {
	baseRegistry := prometheus.NewRegistry()
	// Create a prefixed registry from the base registry
	prefixedRegistry := prometheus.WrapRegistererWithPrefix(namespace+"_", baseRegistry)

	// Register the default Go metrics and process metrics with the prefixed registry
	prefixedRegistry.MustRegister(collectors.NewGoCollector())
	prefixedRegistry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// Custom gauge metric for service health
	serviceHealth := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "service_health",
		Help: "Service health status: 1 for up, 0 for down",
	}, []string{"service"})

	// Register custom metric with the prefixed registry
	prefixedRegistry.MustRegister(serviceHealth)

	return &MetricsServer{
		Registry:      baseRegistry,
		ServiceHealth: serviceHealth,
	}
}

// Handler returns an http.Handler that serves the metrics endpoint
func (m *MetricsServer) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: It is just for simplicity update status to up
		m.ServiceHealth.With(prometheus.Labels{"service": "main"}).Set(1)
		// Use promhttp.HandlerFor with the appropriate registry
		promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	})
}
