package tasks

import (
	"fmt"
	"testing"
)

type testTask struct {
	a int
}

func (this *testTask) Exec() {
	fmt.Println(this.a)
}

func task(goNum int) {
	desc := fmt.Sprintf("task_%v", goNum)
	Register(desc, 100, goNum)

	loop := 10
	for i := 0; i < loop; i++ {
		AddTask(desc, &testTask{i})
	}
}

func taskFunc(goNum int) {
	desc := fmt.Sprintf("taskFunc_%v", goNum)
	Register(desc, 100, goNum)

	loop := 10
	for i := 0; i < loop; i++ {
		k := i
		AddTaskFunc(desc, func() {
			fmt.Println(desc, k)
		})
	}
}

//
func TestAddTask(t *testing.T) {
	task(1)
	Stop()
}

//
func TestAddTaskFunc(t *testing.T) {
	taskFunc(1)
	Stop()
}

//
func TestMulti(t *testing.T) {
	task(2)
	taskFunc(3)
	Stop()
}
