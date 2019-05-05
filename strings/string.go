package strings

import (
	"strconv"
	"strings"

	"util/logs"
)

////////////////////////////////////////////////////////////////////////////////
//
func Split(s, sep string, funcPre func(src string) string) []string {
	//
	if funcPre != nil && s != "" {
		s = funcPre(s)
	}

	if "" == s {
		return nil
	}

	return strings.Split(s, sep)
}

// 如果有错误值，执行continue，并打印日志
func StringArrayToInt64Array(array []string) []int64 {
	var result []int64 = make([]int64, len(array))
	for i, item := range array {
		str, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			logs.Info("string Array To int64 Array，Error!! value=%v, array=%v", item, array)
			continue
		}
		result[i] = str
	}
	return result[:]
}

func Int64ArrayToStringArray(array []int64) []string {
	var result []string = make([]string, len(array))
	for i, item := range array {
		int64Str := strconv.FormatInt(item, 10)
		result[i] = int64Str
	}
	return result[:]
}

//
func ToInt32Array(str, sep string, num int) ([]int32, bool) {
	ss1 := strings.Split(str, sep)
	if len(ss1) != num {
		logs.Error("ToInt32Array|invalid num! str=%v, needNum=%v", str, num)
		return nil, false
	}

	ret := make([]int32, num)
	for i, v := range ss1 {
		vv, e := strconv.Atoi(v)
		if e != nil {
			logs.Error("ToInt32Array|invalid value! str=%v, error=%v", str, e)
			return nil, false
		}
		ret[i] = int32(vv)
	}

	return ret, true
}
