package main

import (
	"github.com/weedge/pipeline-go/pkg/logger"
)

func main() {
	logger.InitLoggerWithConfig(logger.NewDefaultLoggerConfig())

	// Example usage through a helper function
	logInfo := func(msg string, args ...any) {
		logger.Info(msg, args...)
	}

	logInfo("This log should show the caller from main", "key", "value")
}
