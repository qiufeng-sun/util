package tasks

import (
	"sync"
	"util/logs"
)

//
var _ = logs.Debug

// interface
type ITask interface {
	Exec()
}

//
type Tasks struct {
	chTask  chan ITask
	maxTask int
	goNum   int

	desc string
	wg   *sync.WaitGroup
}

func NewTasks(desc string, maxTask, goNum int) *Tasks {
	t := &Tasks{
		chTask:  make(chan ITask, maxTask),
		maxTask: maxTask,
		goNum:   goNum,
		desc:    desc,
		wg:      &sync.WaitGroup{},
	}

	logs.Info("task<%v> start! goroutine num:%v", t.desc, t.goNum)
	for i := 0; i < t.goNum; i++ {
		t.wg.Add(1)
		go t.run(i)
	}

	return t
}

func (this *Tasks) run(index int) {
	logs.Info("task<%v> %v start.", this.desc, index)
	defer logs.Info("task<%v> %v stop.", this.desc, index)

	for ch := range this.chTask {
		ch.Exec()
	}

	this.wg.Done()
}

func (this *Tasks) Stop() {
	close(this.chTask)
	this.wg.Wait()

	logs.Info("task<%v> stop! goroutine num:%v", this.desc, this.goNum)
}

func (this *Tasks) Add(t ITask) bool {
	select {
	case this.chTask <- t:
		return true
	default:
		logs.Warn("task<%v> too many task! discard.", this.desc)
		return false
	}
}

type TaskFunc func()

func (f TaskFunc) Exec() {
	f()
}

func (this *Tasks) AddFunc(f func()) bool {
	return this.Add(TaskFunc(f))
}

//
var g_task = make(map[string]*Tasks) // desc=>*Tasks

func Register(desc string, maxTask, goNum int) {
	_, ok := g_task[desc]
	if ok {
		panic("task exist:" + desc)
	}

	t := NewTasks(desc, maxTask, goNum)
	g_task[desc] = t
}

func Stop() {
	for _, t := range g_task {
		t.Stop()
	}
	g_task = make(map[string]*Tasks)
}

func AddTask(desc string, t ITask) bool {
	m, ok := g_task[desc]
	if !ok {
		panic("task: not register! task:" + desc)
	}

	return m.Add(t)
}

func AddTaskFunc(desc string, f func()) bool {
	return AddTask(desc, TaskFunc(f))
}
