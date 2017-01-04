package safe

import (
	"fmt"
	"sync"

	"util"
)

//
type SafeArray struct {
	array []interface{}
	size  int
	mu    *sync.Mutex
}

// new
//func NewArray(arr []interface{}) *SafeArray {
//	return &SafeArray{array:arr, mu:&sync.Mutex{}, size:len(arr)}
//}
func NewArray(size int) *SafeArray {
	return &SafeArray{array: make([]interface{}, size), mu: &sync.Mutex{}, size: size}
}

// check index
func (a SafeArray) checkIndex(index int) {
	if index >= a.size {
		msg := fmt.Sprintf("index(%v) >= len(%v), caller:%v\n", index, a.size, util.Caller(2))
		panic(msg)
	}
}

// get
func (a SafeArray) Get(index int) interface{} {
	a.checkIndex(index)

	a.mu.Lock()
	defer a.mu.Unlock()

	return a.array[index]
}

// set
func (a *SafeArray) Set(index int, val interface{}) {
	a.checkIndex(index)

	a.mu.Lock()
	a.array[index] = val
	a.mu.Unlock()
}

// size
func (a SafeArray) Size() int {
	return a.size
}
