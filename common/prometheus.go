package common

import (
	"runtime"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMetrics 包含了用于记录请求计数的 Prometheus 指标
type PrometheusMetrics struct {
	requestCounter   *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	requestSizeBytes prometheus.Summary
	memAlloc         prometheus.Gauge
	cpuUsage         prometheus.Gauge
	mu               sync.Mutex // 用于保证并发安全性
}

// NewPrometheusMetrics 创建一个新的 PrometheusMetrics 实例
func NewPrometheusMetrics() *PrometheusMetrics {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: []float64{0.1, 0.3, 1, 3, 10},
		},
		[]string{"method", "endpoint"},
	)

	requestSizeBytes := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "http_request_size_bytes",
			Help: "HTTP request size in bytes.",
		},
	)

	memAlloc := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "mem_alloc_bytes",
			Help: "Total memory allocated by Go runtime, in bytes.",
		},
	)

	cpuUsage := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "CPU usage of the Go process, as a percentage.",
		},
	)

	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestSizeBytes)
	prometheus.MustRegister(memAlloc)
	prometheus.MustRegister(cpuUsage)

	return &PrometheusMetrics{
		requestCounter:   requestCounter,
		requestDuration:  requestDuration,
		requestSizeBytes: requestSizeBytes,
		memAlloc:         memAlloc,
		cpuUsage:         cpuUsage,
	}
}

// IncRequestCounter 增加请求计数
func (p *PrometheusMetrics) IncRequestCounter(method, endpoint string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.requestCounter.WithLabelValues(method, endpoint).Inc()
}

// RecordRequestDuration 记录请求持续时间
func (p *PrometheusMetrics) RecordRequestDuration(method, endpoint string, duration float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.requestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// RecordRequestSize 记录请求大小
func (p *PrometheusMetrics) RecordRequestSize(size float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.requestSizeBytes.Observe(size)
}

// RecordMemoryUsage 记录内存使用情况
func (p *PrometheusMetrics) RecordMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.memAlloc.Set(float64(m.Alloc))
}

// RecordCPUUsage 记录 CPU 使用率
func (p *PrometheusMetrics) RecordCPUUsage() {
	cpuUsage := float64(runtime.NumCPU()) / float64(runtime.NumCPU())
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cpuUsage.Set(cpuUsage)
}

// PrometheusMiddleware 用于记录请求计数和持续时间
func PrometheusMiddleware(metrics *PrometheusMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 记录请求大小
		size := float64(c.Writer.Size())
		metrics.RecordRequestSize(size)

		c.Next()

		duration := time.Since(start).Seconds()

		// 记录请求计数和持续时间
		metrics.IncRequestCounter(c.Request.Method, c.Request.URL.Path)
		metrics.RecordRequestDuration(c.Request.Method, c.Request.URL.Path, duration)
	}
}
