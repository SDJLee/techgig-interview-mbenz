package serve

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SDJLee/mercedes-benz/handler"
	"github.com/SDJLee/mercedes-benz/metrics"
	"github.com/gin-gonic/gin"
)

func Serve(port string) {
	router := setupRouter()
	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Server started and listening on the port: ", port)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("failed to start devero-services-evv", err)
		panic(err)
	}
}

func setupRouter() http.Handler {
	router := gin.New()

	router.Use(metrics.MeasureApiComputationTime())

	apiRoute := router.Group("/api")
	apiRoute.GET("health", handler.HandleHealthCheck)

	apiRouteV1 := apiRoute.Group("/v1")
	apiRouteV1.POST("/check-fuel", handler.HandleFuelCheck)
	return router
}
