package helmfn

import (
	"log"
	"os"
	"strings"
)

var (
	// debugEnabled is set once at package initialization based on LOG_LEVEL env var
	debugEnabled bool
)

func init() {
	// Check LOG_LEVEL environment variable once at startup
	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	debugEnabled = logLevel == "debug"
}

// DebugLog prints debug messages when LOG_LEVEL=debug
func DebugLog(format string, args ...interface{}) {
	if debugEnabled {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// IsDebugEnabled returns true if debug logging is enabled
func IsDebugEnabled() bool {
	return debugEnabled
}