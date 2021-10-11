package handler

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/SDJLee/mercedes-benz/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var requests int64

func HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "breathing...",
	})
}

func HandleFuelCheck(c *gin.Context) {
	var reqBody model.Request
	if err := c.ShouldBindBodyWith(&reqBody, binding.JSON); err != nil {
		fmt.Println("invalid request", err)
		c.String(http.StatusBadRequest, `invalid request`)
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

// TODO: stations in response should have
