/**
Create by SunXiguang 2019-01-09
Desc
*/
package convs

import (
	"testing"
)

//
func TestToInterfaceSlice(t *testing.T) {
	ss := []string{"a", "ax", "99", "iix"}
	is := []int{1, 2, 10, 8, 99}

	r := ToInterfaceSlice(ss[0:1])
	t.Log(r)

	r = ToInterfaceSlice(ss)
	t.Log(r)

	r = ToInterfaceSlice(is[0:1])
	t.Log(r)

	r = ToInterfaceSlice(is)
	t.Log(r)
}
