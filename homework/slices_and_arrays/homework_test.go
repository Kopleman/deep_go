package main

import (
	"reflect"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue struct {
	values []int
	front  int
	rear   int
	size   int
}

func NewCircularQueue(size int) CircularQueue {
	return CircularQueue{
		values: slices.Repeat([]int{-1}, size),
		front:  -1,
		rear:   -1,
		size:   size,
	}
}

func (q *CircularQueue) Push(value int) bool {
	if q.Full() {
		return false
	}
	//first value come in
	if q.front == -1 {
		q.front = 0
		q.rear = 0
	} else {
		q.rear = (q.rear + 1) % q.size
	}

	q.values[q.rear] = value
	return true
}

func (q *CircularQueue) Pop() bool {
	if q.front == -1 {
		return false
	}

	if q.front == q.rear { // edge case poping last element from queue
		front := q.front
		q.front = -1
		q.rear = -1
		q.values[front] = -1
		return true
	}

	q.values[q.front] = -1
	q.front = (q.front + 1) % q.size

	return true
}

func (q *CircularQueue) Front() int {
	if q.Empty() {
		return -1
	}
	return q.values[q.front]
}

func (q *CircularQueue) Back() int {
	if q.Empty() {
		return -1
	}
	return q.values[q.rear]
}

func (q *CircularQueue) Empty() bool {
	return q.front == -1
}

func (q *CircularQueue) Full() bool {
	return (q.rear+1)%q.size == q.front
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

func TestCircularQueue_CyclePushPop(t *testing.T) {
	queue := NewCircularQueue(3)
	assert.True(t, queue.Push(10))
	assert.True(t, queue.Push(20))
	assert.True(t, queue.Push(30))

	assert.True(t, reflect.DeepEqual([]int{10, 20, 30}, queue.values))

	assert.False(t, queue.Push(40))

	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{-1, 20, 30}, queue.values))

	assert.True(t, queue.Push(40))

	assert.True(t, reflect.DeepEqual([]int{40, 20, 30}, queue.values))

	assert.Equal(t, 20, queue.Front())
	assert.Equal(t, 40, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{-1, -1, -1}, queue.values))
	assert.True(t, queue.Empty())
}

func TestCircularQueue_SizeOne(t *testing.T) {
	queue := NewCircularQueue(1)
	assert.True(t, queue.Empty())
	assert.True(t, queue.Push(5))
	assert.True(t, queue.Full())

	assert.True(t, reflect.DeepEqual([]int{5}, queue.values))

	assert.False(t, queue.Push(6))
	assert.Equal(t, 5, queue.Front())
	assert.Equal(t, 5, queue.Back())
	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{-1}, queue.values))

	assert.True(t, queue.Empty())
	assert.True(t, queue.Push(7))
	assert.Equal(t, 7, queue.Front())

	assert.True(t, reflect.DeepEqual([]int{7}, queue.values))

	assert.True(t, queue.Pop())
	assert.True(t, queue.Empty())

	assert.True(t, reflect.DeepEqual([]int{-1}, queue.values))
}

func TestCircularQueue_ManyOperations(t *testing.T) {
	queue := NewCircularQueue(4)
	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))

	assert.True(t, reflect.DeepEqual([]int{1, 2, -1, -1}, queue.values))

	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{-1, 2, -1, -1}, queue.values))

	assert.True(t, queue.Push(3))
	assert.True(t, queue.Push(4))
	assert.True(t, queue.Push(5))

	assert.True(t, reflect.DeepEqual([]int{5, 2, 3, 4}, queue.values))

	assert.False(t, queue.Push(6))
	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 5, queue.Back())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{-1, -1, -1, -1}, queue.values))

	assert.True(t, queue.Empty())
}

func TestCircularQueue_FullEmptyCycle(t *testing.T) {
	queue := NewCircularQueue(2)
	for i := 0; i < 10; i++ {
		assert.True(t, queue.Push(i))
		assert.True(t, queue.Push(i+1))
		assert.True(t, queue.Full())
		assert.Equal(t, i+1, queue.Back())
		assert.True(t, queue.Pop())
		assert.True(t, queue.Pop())
		assert.True(t, queue.Empty())
	}
}

func TestCircularQueue_ValuesNotLost(t *testing.T) {
	queue := NewCircularQueue(3)
	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))

	assert.True(t, reflect.DeepEqual([]int{1, 2, -1}, queue.values))

	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{-1, 2, -1}, queue.values))

	assert.True(t, queue.Push(3))
	assert.True(t, queue.Push(4))
	assert.False(t, queue.Push(5))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{4, -1, -1}, queue.values))

	assert.False(t, queue.Empty())
	assert.Equal(t, 4, queue.Front())
	assert.Equal(t, 4, queue.Back())
}
