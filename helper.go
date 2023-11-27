package decorator

import (
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
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
