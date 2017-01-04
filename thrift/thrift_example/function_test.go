package gamerec

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

//
var (
	// params
	imei        = "nVEzBAuV7GelBWUi"
	count int32 = 10
)

//
func toJsonStr(obj interface{}) string {
	b, e := json.Marshal(obj)
	if e != nil {
		fmt.Printf("error:%v, str:%+v", e, obj)
		return ""
	}

	return string(b)
}

//
func TestNextGames(t *testing.T) {
	// init
	ServiceTimeout = time.Second * 1000
	InitPool("staging", 5, 10)

	var (
		res interface{}
		e   error
	)
	//	//
	//	res, e = NextGames("", count, false, false)
	//	t.Logf("imei=null, res=%+v, error=%v", toJsonStr(res), e)

	//	//
	//	for i := 0; i < 10; i++ {
	//		res, e := NextGames("", count, false, false)
	//		t.Logf("%v: imei=%v, res=%v, error=%v", i, imei, toJsonStr(res), e)
	//	}

	//
	res, e = NextGames("", count, true)
	t.Logf("reset: imei=%v, res=%+v, error=%v", imei, toJsonStr(res), e)

	res, e = NextGames("", count, true)
	t.Logf("reset md5: imei=%v, res=%+v, error=%v", imei, toJsonStr(res), e)
}
