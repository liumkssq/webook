package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	Namespace  string
	Subsystem  string
	Name       string
	Help       string
	InstanceID string
}

func (m *MiddlewareBuilder) Builder() gin.HandlerFunc {
	labels := []string{"method", "pattern", "status"}
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name + "_resp_time",
		Help:      m.Help,
		ConstLabels: map[string]string{
			"instance_id": m.InstanceID,
		},
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
			1:    0,
		},
	}, labels)
	prometheus.MustRegister(summary)
	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			pattern := ctx.FullPath()
			if pattern == "" {
				pattern = "unknown"
			}
			summary.WithLabelValues(ctx.Request.Method, pattern, strconv.Itoa(ctx.Writer.Status())).
				Observe(float64(duration.Milliseconds()))
		}()
		ctx.Next()
	}
}
