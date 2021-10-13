package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SDJLee/mercedes-benz/handler"
	log "github.com/SDJLee/mercedes-benz/logger"
	"github.com/SDJLee/mercedes-benz/util"
	"github.com/spf13/viper"
)

var logger = log.Logger()

func serve() {
	port := viper.GetInt(util.Port)
	writeTimeout := viper.GetInt(util.ServerWriteTimeout)
	readTimeout := viper.GetInt(util.ServerReadTimeout)

	if writeTimeout == 0 {
		writeTimeout = 15
	}
	if readTimeout == 0 {
		readTimeout = 15
	}
	logger.Infof("attempting to serve in port '%d' \n", port)
	fmt.Printf("attempting to serve in port '%d' \n", port)
	router := handler.SetupRouter()
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("localhost:%d", port),
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}
	logger.Info("Server started and listening on the port: ", port)
	fmt.Println("Server started and listening on the port: ", port)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start merc-benz-route-checker", err)
		fmt.Println("failed to start merc-benz-route-checker", err)
		panic(err)
	}
}
