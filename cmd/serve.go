package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SDJLee/mercedes-benz/handler"
	"github.com/SDJLee/mercedes-benz/metrics"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func serve() {
	port := viper.GetInt("PORT")
	writeTimeout := viper.GetInt("SERVER_WRITE_TIMEOUT")
	readTimeout := viper.GetInt("SERVER_READ_TIMEOUT")

	if writeTimeout == 0 {
		writeTimeout = 15
	}
	if readTimeout == 0 {
		readTimeout = 15
	}
	fmt.Printf("attempting to serve in port '%d' \n", port)
	router := setupRouter()
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("localhost:%d", port),
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}
	fmt.Println("Server started and listening on the port: ", port)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("failed to start benz", err)
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
