package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamsainadh/neovault/shared/metrics"
)

func MetricMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		duration := time.Since(start).Seconds()

		statusClass := fmt.Sprintf("%dxx", ctx.Writer.Status()/100)

		metrics.HttpRequestsTotal.WithLabelValues(
			ctx.Request.Method,
			ctx.Request.URL.Path,
			statusClass,
		).Inc()

		metrics.RequestDuration.WithLabelValues(
			ctx.Request.Method,
			ctx.FullPath(),
		).Observe(duration)
	}
}
