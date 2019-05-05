package convs

import (
	"reflect"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
//
func ToInt(v interface{}) int {
	switch v.(type) {
	case int:
		return v.(int)
	case float64:
		return int(v.(float64))
	case string:
		d, _ := strconv.Atoi(v.(string))
		return d
	case int32:
		return int(v.(int32))
	case float32:
		return int(v.(float32))
	case int64:
		return int(v.(int64))
	}

	return 0
}

//
func ToFloat64(v interface{}) float64 {
	switch v.(type) {
	case float64:
		return v.(float64)
	case string:
		r, _ := strconv.ParseFloat(v.(string), 64)
		return r
	}

	return 0.0
}

//
func ToInterfaceSlice(val interface{}) []interface{} {
	rv := reflect.ValueOf(val)
	num := rv.Len()
	ret := make([]interface{}, num)
	for i := 0; i < num; i++ {
		ret[i] = rv.Index(i).Interface()
	}

	return ret
}
