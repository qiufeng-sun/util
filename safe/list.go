// 协程安全的list(对contain/list进行封装)
package safe

//
import (
	"container/list"
	"sync"
)

//
type SafeList struct {
	lst *list.List
	mu  *sync.Mutex
}

// new
func NewList() *SafeList {
	return &SafeList{list.New(), &sync.Mutex{}}
}

// 遍历接口
//	el := l.Begin()
//	for ; el != nil; el = el.Next() {
//		// do something
//	}
//	l.End()
func (l *SafeList) Begin() *list.Element {
	l.mu.Lock()

	return l.lst.Front()
}

func (l *SafeList) End() {
	l.mu.Unlock()
}

//
func (l *SafeList) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.lst.Len()
}

//
func (l *SafeList) Clear() {
	l.mu.Lock()
	l.lst.Init()
	l.mu.Unlock()
}

// push back
func (l *SafeList) PushBack(v interface{}) *list.Element {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.lst.PushBack(v)
}

//
