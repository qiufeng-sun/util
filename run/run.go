package run

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"

	"util/logs"
)

// 调用者
func Caller(steps int) string {
	if pc, _, line, ok := runtime.Caller(steps + 1); ok {
		return fmt.Sprintf("[%s:%d]", runtime.FuncForPC(pc).Name(), line)
	}

	return "[?]"
}

func CallerFile(steps int) string {
	if _, filename, line, ok := runtime.Caller(steps + 1); ok {
		return fmt.Sprintf("[%s:%d]", filename, line)
	}

	return "[?]"
}

// CheckTime: 记录进入和退出函数时间, 并根据条件输出.
//   例:在需要统计的函数起始处调用defer CheckTime(200)()
func CheckTime(outMilliSec int64, steps int, debug bool) func() {
	start := time.Now()
	return func() {
		end := time.Now()
		t := (int64)(end.Sub(start) / 1000000)
		if debug || t > outMilliSec {
			if steps < 0 {
				steps = 0
			}
			caller := Caller(steps + 1)
			logs.Warn("caution%v: start:%v, end:%v, elapsed:%v", caller, start, end, t)
		}
	}
}

// usage: 在goroutine开始时执行 defer PrintPanic()
func PrintPanic(exit bool) {
	if r := recover(); r != nil {
		logs.Critical("panic:%v", r)
		logs.Critical("%s", debug.Stack())

		if exit {
			logs.Critical("exit now!")
			os.Exit(1)
		}
	}
	logs.GetLogger().Flush()
}

// usage: Exec(func, param1, param2, ...)
func Exec(exit bool, f interface{}, params ...interface{}) {
	vf := reflect.ValueOf(f)
	vps := make([]reflect.Value, len(params))
	for i := 0; i < len(params); i++ {
		vps[i] = reflect.ValueOf(params[i])
	}

	defer PrintPanic(exit)
	vf.Call(vps)
}
