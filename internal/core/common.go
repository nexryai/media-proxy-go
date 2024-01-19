package core

import (
	"os"
)

func IsDebugMode() bool {
	return os.Getenv("DEBUG") == "1"
}
