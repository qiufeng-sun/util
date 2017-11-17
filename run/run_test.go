package run

import (
	"fmt"
	"testing"
	"time"
)

//
func TestRecover(t *testing.T) {
	go func() {
		defer Recover(false)

		panic("hehe!")
	}()

	time.Sleep(time.Second)
}

//
func print0() {
	fmt.Println("no param function")
	panic("f0")
}

func print1(s string) {
	fmt.Println(s)
	panic("f1")
}

func print2(s, s2 string) {
	fmt.Println(s, s2)
	panic("f2")
}

//
func TestGoExec(t *testing.T) {
	go Exec(false, print0)
	go Exec(false, print1, "f1")
	go Exec(false, print2, "f2", "f22")

	go func() {
		defer Recover(false)
		print0()
	}()

	time.Sleep(time.Second)
}
