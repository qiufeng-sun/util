package logs

import (
	"testing"
	"time"
)

//
const cfgFile = "log.conf"

//
func TestLog(t *testing.T) {
	Init(cfgFile)

	Debug("debug lv log!")
	Info("info lv log!")
	Warn("warn")

	Debugln("debugln")
}

//
func TestLogAsync(t *testing.T) {
	Init(cfgFile).Async(0)

	Debug("debug lv log!")
	Info("info lv log!")
	Warn("warn")

	// wait log write
	time.Sleep(time.Second)
}

//
func TestCheckTime(t *testing.T) {
	f := func(milliSec int64) {
		defer CheckTime(milliSec, 0, false)()
		time.Sleep(time.Millisecond * 200)
	}

	go f(100)
	go f(300)

	time.Sleep(time.Second)
}
