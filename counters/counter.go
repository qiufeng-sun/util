// 计数器: 调用自增接口时, 如果与记录中key相同则增加记录中计数, 否则重新计数
// 例: 限制每分钟(秒钟)接口访问次数, 以时间为参数调用自增接口, 根据返回值判断是否到上限
package counters

import "sync"

////////////////////////////////////////////////////////////////////////////////
//
type TimesCounter struct {
	key string // 20150914
	num int
}

//
type AtomicTimesCounter struct {
	*sync.Mutex
	*TimesCounter
}

//
func NewAtomicTimesCounter() *AtomicTimesCounter {
	return &AtomicTimesCounter{
		Mutex:        &sync.Mutex{},
		TimesCounter: &TimesCounter{},
	}
}

//
func (this *AtomicTimesCounter) GetTimes(key string) int {
	this.Lock()
	defer this.Unlock()

	d := this.TimesCounter
	if d.key != key {
		return 0
	}

	return d.num
}

//
func (this *AtomicTimesCounter) IncTimes(key string) int {
	this.Lock()
	defer this.Unlock()

	d := this.TimesCounter
	if d.key != key {
		d.key = key
		d.num = 1
	} else {
		d.num++
	}

	return d.num
}
