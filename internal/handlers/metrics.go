package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	InferenceRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "inference_requestes_total",
			Help: "Total number of inference requests",
		},
		[]string{"models", "version", "status"},
	)

	InferenceLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "inference_;atency_ms",
			Help:    "Inference latency in milliseconds",
			Buckets: []float64{10, 50, 100, 250, 500, 1000, 2000},
		},
		[]string{"model", "version"},
	)
)

func init() {
	prometheus.MustRegister(InferenceRequests)
	prometheus.MustRegister(InferenceLatency)
}

func Metrices() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP{c.Writer, c.Request}
	}
}
