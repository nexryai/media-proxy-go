package logger

import (
	"github.com/nexryai/visualog"
	"os"
)

func GetLogger(moduleName string) *visualog.Logger {
	var debugMode bool
	if os.Getenv("DEBUG") == "1" {
		debugMode = true
	}

	return &visualog.Logger{
		ModuleName: moduleName,
		ShowDebug:  debugMode,
		ShowTime:   true,
		ShowCaller: true,
		ShowTrace:  true,
	}
}
