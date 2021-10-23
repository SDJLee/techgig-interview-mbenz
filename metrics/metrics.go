package metrics

import (
	"time"

	log "github.com/SDJLee/mercedes-benz/logger"
	"github.com/gin-gonic/gin"
)

// The idea of this package is to monitor and collect metrics for the service for instrumentation.
// Current implementation is a scaffold and the instrumentation are logged. This should be enhanced to connect
// with an instrumentation tool.

var metricsLogger = log.SubLogger("metrics")

// MonitorTimeElapsed reports the time elapsed to execute a function.
// To use this, add it in defer at beginning of a function like below code.
// defer helper.MonitorTimeElapsed("task-name")()
func MonitorTimeElapsed(what string) func() {
	start := time.Now()
	return func() {
		metricsLogger.Infof("%s took %v\n", what, time.Since(start))
	}
}

// MeasureApiComputationTime is a middleware function that
// reports the time elapsed for an API to execute.
func MeasureApiComputationTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			metricsLogger.Infof("%s took %v\n", c.Request.URL.Path, time.Since(start))
		}()
		c.Next()
	}
}
