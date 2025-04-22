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

	"github.com/vmihailenco/msgpack/v5"
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

func marshalMsgpack(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

// UnmarshalMsgpack 解析 MessagePack 格式的字节切片 data 并将结果存储在 v 指向的值中。
// v 必须是一个指向目标数据结构（如结构体、map、slice、基本类型等）的指针，类似于 encoding/json.Unmarshal。
// 如果 v 是 nil 或者不是指针，UnmarshalMsgpack 会返回错误。
//
// 对于数字类型：
//   - 当反序列化到 interface{} 或 map[string]any 时，msgpack 通常会尝试将整数
//     解码为 int64 或 uint64，浮点数解码为 float64，这比 json 默认将所有数字
//     解码为 float64 更能保留类型信息。
//   - 当反序列化到具体的类型（如 *int, *float32, *MyStruct）时，会进行相应的类型转换。
//
// 对于结构体：
//   - 默认情况下，它会匹配导出字段名（大小写敏感）或 'msgpack' 标签。
//   - 嵌套的结构体如果存在于 map[string]any 的值中，反序列化到 map[string]any 时
//     通常会被解码为 map[string]interface{}。
func unmarshalMsgpack(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}
