package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"
)

//////////////////////////////////////////////////////////////////////////////
//
const (
	X_MinuteSec   int64 = 60
	X_HourSec     int64 = 60 * 60         // hour sec
	X_DaySec      int64 = X_HourSec * 24  // day sec
	X_DayMilliSec int64 = X_DaySec * 1000 // day millisec

	X_TimeZone_ModSec int64 = X_HourSec * 8 // 8 hours
)

var (
	e_DuplicateSrc_Nil = errors.New("copy source is nil")
	e_Bytes_Invalid    = errors.New("byte slice invalid")
)

//////////////////////////////////////////////////////////////////////////////
//
func CurMillisecond() int64 {
	return time.Now().UnixNano() / 1000000
}

func ToJsonString(v interface{}) string {
	b, e := json.Marshal(v)
	if e != nil {
		return "ToString() failed!"
	}

	return string(b)
}

func ToJsonBytes(v interface{}) []byte {
	b, e := json.Marshal(v)
	if e != nil {
		return []byte("ToString() failed!")
	}

	return b
}

// 只能复制Public字段
func Duplicate(src interface{}, dst interface{}) error {
	//
	if nil == src {
		return e_DuplicateSrc_Nil
	}

	//
	b, e := json.Marshal(src)
	if e != nil {
		return e
	}

	//
	if e := json.Unmarshal(b, dst); e != nil {
		return e
	}

	return nil
}

//
func ToObj(b []byte, dst interface{}) error {
	//
	if nil == b || len(b) < 2 {
		return e_Bytes_Invalid
	}

	//
	if e := json.Unmarshal(b, dst); e != nil {
		return e
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
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

////////////////////////////////////////////////////////////////////////////////
//
func ToValues(vs ...interface{}) []reflect.Value {
	if len(vs) <= 0 {
		return nil
	}

	r := make([]reflect.Value, len(vs))
	for i, v := range vs {
		r[i] = reflect.ValueOf(v)
	}

	return r
}
