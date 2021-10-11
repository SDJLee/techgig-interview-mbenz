package metrics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// MonitorTimeElapsed reports the time elapsed to execute a function.
// To use this, add it in defer at beginning of a function like below code.
// defer helper.MonitorTimeElapsed("task-name")()
func MonitorTimeElapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

// MeasureApiComputationTime is a middleware function that
// reports the time elapsed for an API to execute.
func MeasureApiComputationTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			fmt.Printf("%s took %v\n", c.Request.URL.Path, time.Since(start))
		}()
		c.Next()
	}
}
