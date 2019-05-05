package util

import (
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
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

// DrainChannel waits for the channel to finish
// emptying (draining) for up to the expiration.  It returns
// true if the drain completed (the channel is empty), false otherwise.
func DrainChannel(ch reflect.Value, expire time.Time) bool {
	var dur = time.Millisecond * 10

	for {
		if ch.Len() == 0 {
			return true
		}
		now := time.Now()
		if now.After(expire) {
			return false
		}
		// We sleep the lesser of the remaining time, or
		// 10 milliseconds.  This polling is kind of suboptimal for
		// draining, but its far far less complicated than trying to
		// arrange special messages to force notification, etc.
		//dur = expire.Sub(now)
		//if dur > time.Millisecond*10 {
		//	dur = time.Millisecond * 10
		//}
		time.Sleep(dur)
	}
}

// 参考: rand.Perm
func Perm(n int32) []int32 {
	m := make([]int32, n)

	for i := int32(0); i < n; i++ {
		j := rand.Int31n(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}

//把数组打乱顺序
func PermArray(array []int32) []int32 {
	for i := int32(0); i < int32(len(array)); i++ {
		j := rand.Int31n(i + 1)
		temp := array[i]
		array[i] = array[j]
		array[j] = temp
	}
	return array
}

func ProcTime(t string) (int64, bool) {
	r, e := time.Parse("2006010215", t)
	if e != nil {
		return 0, false
	}

	return r.Unix(), true
}
