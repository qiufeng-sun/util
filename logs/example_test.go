package logs

import (
	"testing"
	"time"

	ghLogs "github.com/astaxie/beego/logs"
)

//
func TestFile(t *testing.T) {
	log := ghLogs.NewLogger(10000)
	// use 0666 as test perm cause the default umask is 022
	log.SetLogger("file", `{"filename":"test.log", "perm": "0666"}`)
	log.Debug("debug")
	log.Info("info")
}

//
func TestFileAsync(t *testing.T) {
	log := ghLogs.NewLogger(10000).Async(0)
	// use 0666 as test perm cause the default umask is 022
	log.SetLogger("file", `{"filename":"test.log", "perm": "0666"}`)
	log.Warn("warn")

	time.Sleep(time.Second)
}
