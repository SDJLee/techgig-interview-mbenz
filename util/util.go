package util

import (
	"fmt"
	"os"
)

func GetEnv() string {
	env := os.Getenv(AppEnv)
	fmt.Printf("getEnv from os '%v'\n", env)
	if env == "" {
		env = EnvDev
	}
	fmt.Println("getEnv", env)
	return env
}
