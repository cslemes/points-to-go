package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Definindo as métricas Prometheus
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de requisições HTTP recebidas.",
		},
		[]string{"handler", "method"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duração das requisições HTTP em segundos.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler"},
	)
)

func init() {
	// Registrando as métricas no Prometheus
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

func MetricsMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()

	// Atualizando as métricas após a requisição
	handler := c.HandlerName()
	requestCounter.WithLabelValues(handler, c.Request.Method).Inc()
	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues(handler).Observe(duration)
}
