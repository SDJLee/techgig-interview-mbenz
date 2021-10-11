package cmd

import (
	"fmt"
	"os"

	"github.com/SDJLee/mercedes-benz/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	EnvDev    = "dev"
	EnvProd   = "prod"
	appEnvStr = "APP_ENV"
)

var rootCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		// appEnv, _ := cmd.Flags().GetString(appEnvStr)
		// setupAppEnv(appEnv)

		logger.SetupLogger()
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
	// rootCmd.PersistentFlags().StringP(appEnvStr, "e", "", "Development or production environment")
	loadConfig()
}

// func setupAppEnv(appEnv string) {
// 	if appEnv == "" || appEnv == EnvDev {
// 		os.Setenv(appEnvStr, EnvDev)
// 		return
// 	}
// 	os.Setenv(appEnv, EnvProd)
// }

func loadConfig() {
	fmt.Println("env str", os.Getenv(appEnvStr))
	env := os.Getenv(appEnvStr)
	if env == "" {
		env = EnvDev
	}
	viper.AddConfigPath("/app")
	viper.SetConfigName(fmt.Sprintf("app-%s", env))
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return
	}

}
