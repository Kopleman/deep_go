package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type IntNum interface {
	int8 | int16 | int32 | int64 | int
}

type CircularQueue[T IntNum] struct {
	values []T
	front  int
	rear   int
	length int
}

func NewCircularQueue[T IntNum](size int) CircularQueue[T] {
	return CircularQueue[T]{
		values: make([]T, size),
		rear:   -1,
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	q.rear = (q.rear + 1) % len(q.values)
	q.values[q.rear] = value
	q.length++

	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}

	q.front = (q.front + 1) % len(q.values)
	q.length--

	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.front]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}

	v := q.values[q.rear]
	return v
}

func (q *CircularQueue[T]) Empty() bool {
	return q.length == 0
}

func (q *CircularQueue[T]) Full() bool {
	return q.length == len(q.values)
}

func TestCircularQueue_Generic(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		queue := NewCircularQueue[int](3)
		assert.True(t, queue.Push(1))
		assert.True(t, queue.Push(2))
		assert.True(t, queue.Push(3))
		assert.False(t, queue.Push(4))
		assert.Equal(t, 1, queue.Front())
		assert.Equal(t, 3, queue.Back())
		assert.True(t, queue.Pop())
		assert.True(t, queue.Push(4))
		assert.Equal(t, 2, queue.Front())
		assert.Equal(t, 4, queue.Back())
	})
	t.Run("int8", func(t *testing.T) {
		queue := NewCircularQueue[int8](2)
		assert.True(t, queue.Push(10))
		assert.True(t, queue.Push(20))
		assert.False(t, queue.Push(30))
		assert.Equal(t, int8(10), queue.Front())
		assert.Equal(t, int8(20), queue.Back())
		assert.True(t, queue.Pop())
		assert.True(t, queue.Push(30))
		assert.Equal(t, int8(20), queue.Front())
		assert.Equal(t, int8(30), queue.Back())
	})
	t.Run("int16", func(t *testing.T) {
		queue := NewCircularQueue[int16](2)
		assert.True(t, queue.Push(1000))
		assert.True(t, queue.Push(2000))
		assert.False(t, queue.Push(3000))
		assert.Equal(t, int16(1000), queue.Front())
		assert.Equal(t, int16(2000), queue.Back())
		assert.True(t, queue.Pop())
		assert.True(t, queue.Push(3000))
		assert.Equal(t, int16(2000), queue.Front())
		assert.Equal(t, int16(3000), queue.Back())
	})
	t.Run("int32", func(t *testing.T) {
		queue := NewCircularQueue[int32](2)
		assert.True(t, queue.Push(100000))
		assert.True(t, queue.Push(200000))
		assert.False(t, queue.Push(300000))
		assert.Equal(t, int32(100000), queue.Front())
		assert.Equal(t, int32(200000), queue.Back())
		assert.True(t, queue.Pop())
		assert.True(t, queue.Push(300000))
		assert.Equal(t, int32(200000), queue.Front())
		assert.Equal(t, int32(300000), queue.Back())
	})
	t.Run("int64", func(t *testing.T) {
		queue := NewCircularQueue[int64](2)
		assert.True(t, queue.Push(10000000000))
		assert.True(t, queue.Push(20000000000))
		assert.False(t, queue.Push(30000000000))
		assert.Equal(t, int64(10000000000), queue.Front())
		assert.Equal(t, int64(20000000000), queue.Back())
		assert.True(t, queue.Pop())
		assert.True(t, queue.Push(30000000000))
		assert.Equal(t, int64(20000000000), queue.Front())
		assert.Equal(t, int64(30000000000), queue.Back())
	})
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

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
	queue := NewCircularQueue[int](3)
	assert.True(t, queue.Push(10))
	assert.True(t, queue.Push(20))
	assert.True(t, queue.Push(30))

	assert.True(t, reflect.DeepEqual([]int{10, 20, 30}, queue.values))

	assert.False(t, queue.Push(40))

	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{10, 20, 30}, queue.values))

	assert.True(t, queue.Push(40))

	assert.True(t, reflect.DeepEqual([]int{40, 20, 30}, queue.values))

	assert.Equal(t, 20, queue.Front())
	assert.Equal(t, 40, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{40, 20, 30}, queue.values))
	assert.True(t, queue.Empty())
}

func TestCircularQueue_SizeOne(t *testing.T) {
	queue := NewCircularQueue[int](1)
	assert.True(t, queue.Empty())
	assert.True(t, queue.Push(5))
	assert.True(t, queue.Full())

	assert.True(t, reflect.DeepEqual([]int{5}, queue.values))

	assert.False(t, queue.Push(6))
	assert.Equal(t, 5, queue.Front())
	assert.Equal(t, 5, queue.Back())
	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{5}, queue.values))

	assert.True(t, queue.Empty())
	assert.True(t, queue.Push(7))
	assert.Equal(t, 7, queue.Front())

	assert.True(t, reflect.DeepEqual([]int{7}, queue.values))

	assert.True(t, queue.Pop())
	assert.True(t, queue.Empty())

	assert.True(t, reflect.DeepEqual([]int{7}, queue.values))
}

func TestCircularQueue_ManyOperations(t *testing.T) {
	queue := NewCircularQueue[int](4)
	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 0, 0}, queue.values))

	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{1, 2, 0, 0}, queue.values))

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

	assert.True(t, reflect.DeepEqual([]int{5, 2, 3, 4}, queue.values))

	assert.True(t, queue.Empty())
}

func TestCircularQueue_FullEmptyCycle(t *testing.T) {
	queue := NewCircularQueue[int](2)
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
	queue := NewCircularQueue[int](3)
	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 0}, queue.values))

	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{1, 2, 0}, queue.values))

	assert.True(t, queue.Push(3))
	assert.True(t, queue.Push(4))
	assert.False(t, queue.Push(5))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.Equal(t, 4, queue.Front())
	assert.Equal(t, 4, queue.Back())
}
