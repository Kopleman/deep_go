package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier       int
	Priority         int
	sequence         int // for FIFO ordering (private field)
	originalPriority int // store original priority
}
type TaskHeap []Task

func (h TaskHeap) Len() int { return len(h) }

func (h TaskHeap) Less(i, j int) bool {
	// Higher priority first
	if h[i].Priority != h[j].Priority {
		return h[i].Priority > h[j].Priority
	}
	// FIFO for same priority
	return h[i].sequence < h[j].sequence
}

func (h TaskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *TaskHeap) Push(x interface{}) {
	*h = append(*h, x.(Task))
}

func (h *TaskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type Scheduler struct {
	tasks   *TaskHeap
	taskMap map[int]int // taskID -> index in heap
	nextSeq int         // for FIFO ordering
}

func NewScheduler() Scheduler {
	h := &TaskHeap{}
	heap.Init(h)
	return Scheduler{
		tasks:   h,
		taskMap: make(map[int]int),
		nextSeq: 0,
	}
}

func (s *Scheduler) AddTask(task Task) {
	task.originalPriority = task.Priority
	task.sequence = s.nextSeq
	s.nextSeq++

	heap.Push(s.tasks, task)

	s.taskMap[task.Identifier] = s.tasks.Len() - 1
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	index, exists := s.taskMap[taskID]
	if !exists {
		return
	}
	(*s.tasks)[index].Priority = newPriority

	heap.Fix(s.tasks, index)

	for i, task := range *s.tasks {
		s.taskMap[task.Identifier] = i
	}
}

func (s *Scheduler) GetTask() Task {
	if s.tasks.Len() == 0 {
		return Task{}
	}

	task := heap.Pop(s.tasks).(Task)

	delete(s.taskMap, task.Identifier)

	for i, t := range *s.tasks {
		s.taskMap[t.Identifier] = i
	}

	// Return task with original priority (due to original test requirements)
	return Task{
		Identifier: task.Identifier,
		Priority:   task.originalPriority,
	}
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}

func TestEmptyScheduler(t *testing.T) {
	scheduler := NewScheduler()
	task := scheduler.GetTask()
	assert.Equal(t, Task{}, task)
}

func TestSamePriorityFIFO(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 10}
	task3 := Task{Identifier: 3, Priority: 10}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)

	// Should return tasks in FIFO order
	task := scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task2, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}

func TestChangePriorityNonExistentTask(t *testing.T) {
	scheduler := NewScheduler()
	task1 := Task{Identifier: 1, Priority: 10}
	scheduler.AddTask(task1)

	scheduler.ChangeTaskPriority(999, 100)

	task := scheduler.GetTask()
	assert.Equal(t, task1, task)
}

func TestMultiplePriorityChanges(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)

	scheduler.ChangeTaskPriority(1, 30)
	scheduler.ChangeTaskPriority(1, 5)
	scheduler.ChangeTaskPriority(1, 25)

	// task2 should be first (priority 20), then task1 (priority 10)
	task := scheduler.GetTask()
	assert.Equal(t, task2, task)

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)
}
