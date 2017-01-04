package counters

import (
	"sync"
	"testing"
)

//
func TestTimesCounter(t *testing.T) {
	//
	counter := NewAtomicTimesCounter()

	// 0
	key := "20150914"
	num := counter.GetTimes(key)
	if num != 0 {
		t.Errorf("times counter key:%v, num=%v, expected num: 0", key, num)
	}

	// add parallel
	var wg sync.WaitGroup
	total := 10000
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			counter.IncTimes(key)
			wg.Done()
		}()
	}

	// check result
	wg.Wait()
	num = counter.GetTimes(key)
	if num != total {
		t.Errorf("times counter key:%v, num=%v, expected num: %v", key, num, total)
	}

	// another key
	key = "20150915"
	num = counter.GetTimes(key)
	if num != 0 {
		t.Errorf("times counter key:%v, num=%v, expected num: 0", key, num)
	}
}
