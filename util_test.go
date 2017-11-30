package util

import (
	"testing"
	"time"
	"reflect"
)

//
func TestDrainChannel(t *testing.T) {
	ch := make(chan int, 10)
	ch <- 1
	ch <- 2
	go func() {
		time.Sleep(time.Microsecond*100)
		<-ch
		time.Sleep(time.Microsecond*100)
		<-ch
	}()

	var r bool
	r = DrainChannel(reflect.ValueOf(ch), time.Now().Add(time.Second))
	t.Log("channel len", len(ch), r)

	ch <-3
	r = DrainChannel(reflect.ValueOf(ch), time.Now().Add(time.Second))
	t.Log("channel len", len(ch), r)
}
