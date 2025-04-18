package gofnext

import (
	"log/slog"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"
)

var VERSION string = "v0.0.0"

func init() {
	getVersion()
}

func getVersion() string {
	_, filename, _, _ := runtime.Caller(0)
	versionFile := path.Dir(filename) + "/version"
	version, _ := os.ReadFile(versionFile)
	VERSION = strings.TrimSpace(string(version))
	return VERSION

}

func sleepRandom(min time.Duration, max time.Duration) {
	randomDuration := time.Duration(rand.Int63n(int64(max-min)) + int64(min))

	time.Sleep(randomDuration)
}

/*** slogger ************************************/
var slogger *slog.Logger

func init() {
	slogger = getSlogger(false)
}
func getSlogger(isJson bool) (logger *slog.Logger) {
	handlerOpts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
	if isJson {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stderr, handlerOpts))
	}
	// logger2 := logger.With("url", "http://example.com/a/b/c")
	return logger
}

func AssertEqual[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if expected != actual {
		t.Fatalf("Actual %v, expected %v", actual, expected)
	}
}
