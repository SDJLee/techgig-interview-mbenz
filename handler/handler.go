package handler

import (
	"net/http"
	"sync/atomic"

	log "github.com/SDJLee/mercedes-benz/logger"
	"github.com/SDJLee/mercedes-benz/metrics"
	"github.com/SDJLee/mercedes-benz/model"
	"github.com/SDJLee/mercedes-benz/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var requests int64
var logger = log.SubLogger("merc-benz-route-checker")

func HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "breathing...",
	})
}

func HandleFuelCheck(c *gin.Context) {
	var reqBody model.Request
	if err := c.ShouldBindBodyWith(&reqBody, binding.JSON); err != nil {
		logger.Error("invalid request", err)
		c.String(http.StatusBadRequest, `invalid request`)
		return
	}
	incrementRequestCount()
	response := computeArrival(&reqBody, getRequests())
	c.JSON(http.StatusOK, response)
}

func incrementRequestCount() {
	atomic.AddInt64(&requests, 1)
}

func getRequests() int64 {
	return atomic.LoadInt64(&requests)
}

func SetupRouter() http.Handler {
	router := gin.New()

	router.Use(metrics.MeasureApiComputationTime())

	apiRoute := router.Group(util.ApiBasePath)
	apiRoute.GET(util.ApiHealthCheck, HandleHealthCheck)

	apiRouteV1 := apiRoute.Group(util.ApiV1)
	apiRouteV1.POST(util.ApiComputeRoute, HandleFuelCheck)
	return router
}
