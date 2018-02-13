package ut

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdd(t *testing.T) {
	Convey("Add", t, func() {
		So(Add(1, 2), ShouldEqual, 3)
	})
}

func TestEcho(t *testing.T) {
	Convey("Echo", t, func() {
		s := "tessttt"
		So(Echo(s), ShouldEqual, s)
	})
}

func TestAppend(t *testing.T) {
	Convey("Append", t, func() {
		a := []int{1, 2, 3}
		b := []int{2, 3}
		expect := []int{1, 2, 3, 2, 3}
		r := Append(a, b)
		So(r, ShouldNotEqual, expect)
		So(r, ShouldResemble, expect)
	})
}
