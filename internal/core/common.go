package core

import (
	"fmt"
	"os"
	"runtime"
)

func IsDebugMode() bool {
	return os.Getenv("DEBUG_MODE") == "1"
}

func GetProxyConfig() string {
	return os.Getenv("PROXY")
}

func RaisePanicOnHighMemoryUsage(threshold float64) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Calculate memory usage percentage
	usedMemory := float64(memStats.Alloc)
	totalMemory := float64(memStats.Sys)
	memoryUsage := (usedMemory / totalMemory) * 100

	if memoryUsage >= threshold {
		panic(fmt.Errorf("Memory usage exceeded %.2f%% threshold", threshold))
	}
}
