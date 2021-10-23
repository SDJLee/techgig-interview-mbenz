package handler

import (
	"net/http"
	"os"
	"testing"

	"github.com/SDJLee/mercedes-benz/util"
	"github.com/spf13/viper"
)

var router http.Handler

// TestMain for package 'controller' to setup data for testing.
func TestMain(m *testing.M) {
	setup()
	code := m.Run() // invokes the current test function
	cleanup()
	os.Exit(code)
}

// This method sets up router and default viper config.
func setup() {
	router = SetupRouter()
	setupTestConfig()
}

func cleanup() {
	router = nil
}

func setupTestConfig() {
	// testing config
	viper.Set(util.Port, "8080")
	viper.Set(util.AppEnv, util.EnvDev)
	viper.Set(util.ApiAddress, "https://restmock.techgig.com/merc")
}
