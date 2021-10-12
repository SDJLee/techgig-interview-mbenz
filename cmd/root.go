package cmd

import (
	"fmt"
	"os"

	"github.com/SDJLee/mercedes-benz/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the merc-benz-route-checker service",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	loadConfig()
}

func loadConfig() {
	os.Setenv(util.AppEnvStr, util.EnvDev)
	os.Setenv(util.BasePath, "/mnt/d/Developer/go/src/github.com/SDJLee/mercedes-benz")
	fmt.Println("env str", os.Getenv(util.AppEnvStr))
	env := os.Getenv(util.AppEnvStr)
	if env == "" {
		env = util.EnvDev
	}

	fmt.Println("base path str", os.Getenv(util.BasePath))
	basePath := os.Getenv(util.BasePath)
	if basePath == "" {
		basePath = util.DefaultBasePath
	}

	viper.AddConfigPath(basePath)
	viper.SetConfigName(fmt.Sprintf(util.ConfigFileFormat, env))
	viper.SetConfigType(util.ConfigFileType)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return
	}
}
