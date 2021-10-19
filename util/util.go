package util

import (
	"fmt"
	"os"
)

func GetEnv() string {
	env := os.Getenv(AppEnv)
	if env == "" {
		env = EnvDev
	}
	fmt.Println("getEnv", env)
	return env
}
