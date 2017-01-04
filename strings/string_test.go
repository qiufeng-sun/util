package strings

import (
	"strings"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////
//
func TestSplit(t *testing.T) {
	//
	s := "1, 2, testing3"

	ss := Split(s, ",", nil)

	t.Logf("string:%v, slice:%v, num:%v", s, ss, len(ss))

	//
	s = ",,1, 2, testing3,,"

	ss = Split(s, ",", func(src string) string {
		return strings.Trim(src, ",")
	})

	t.Logf("string:%v, slice:%v, num:%v", s, ss, len(ss))

	//
	s = ",,,,"

	ss = Split(s, ",", func(src string) string {
		return strings.Trim(src, ",")
	})

	t.Logf("string:%v, slice:%v, num:%v", s, ss, len(ss))
}
