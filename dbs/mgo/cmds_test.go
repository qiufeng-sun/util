/**
Create by SunXiguang 2019-01-15
Desc
*/
package mgo

import (
	"testing"

	"util"
)

//
func TestMOp(t *testing.T) {
	m := NewM().Op("set", "accId", 981).
		Op("unset", "sex", 1).Op("inc", "age", 1).Op("set", "age", 10)

	t.Log("normal m op:", util.ToJsonString(m.M))

	defer func() {
		if e := recover(); e != nil {
			t.Log("panic:", e, util.ToJsonString(m.M))
		}
	}()

	m.Op("set", "accId", 99)
}
