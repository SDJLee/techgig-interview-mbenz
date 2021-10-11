package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// https://developpaper.com/using-zap-log-library-in-go-language-project-translation/

func SetupLogger() {
	logger, _ := zap.NewProduction()
	fmt.Print(logger)
}
