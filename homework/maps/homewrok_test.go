package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go
type node[K comparable, V any] struct {
	key  K
	data *V
	next *node[K, V]
	prev *node[K, V]
}

func (n *node[K, V]) fetchNode(key K, comparator func(a, b K) bool) *node[K, V] {
	if n.key == key {
		return n
	}
	if comparator(n.key, key) {
		if n.prev == nil {
			n.prev = &node[K, V]{key: key, next: n}
			return n.prev
		}

		if comparator(key, n.prev.key) {
			newNode := &node[K, V]{key: key, prev: n.prev, next: n}
			n.prev.next = newNode
			n.prev = newNode
			return n.prev
		}

		return n.prev.fetchNode(key, comparator)
	}

	if n.next == nil {
		n.next = &node[K, V]{key: key, prev: n}
		return n.next
	}

	if comparator(n.next.key, key) {
		newNode := &node[K, V]{key: key, prev: n, next: n.next}
		n.next.prev = newNode
		n.next = newNode
		return n.next
	}

	return n.next.fetchNode(key, comparator)
}

func (n *node[K, V]) findNode(key K, comparator func(a, b K) bool) *node[K, V] {
	if n.key == key {
		return n
	}
	if comparator(n.key, key) {
		if n.prev == nil {
			return nil
		}

		if comparator(key, n.prev.key) {
			return nil
		}

		return n.prev.findNode(key, comparator)
	}

	if n.next == nil {
		return nil
	}

	if comparator(n.next.key, key) {
		return nil
	}

	return n.next.findNode(key, comparator)
}

func (n *node[K, V]) first() *node[K, V] {
	if n.prev == nil {
		return n
	}

	return n.prev.first()
}

func (n *node[K, V]) last() *node[K, V] {
	if n.next == nil {
		return n
	}

	return n.next.last()
}

type OrderedMap[K comparable, V any] struct {
	root       *node[K, V]
	size       int
	comparator func(a, b K) bool
}

func NewOrderedMap[K comparable, V any](comparator func(a, b K) bool) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{comparator: comparator}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	if m.root == nil {
		m.root = &node[K, V]{
			data: &value,
			key:  key,
		}
		m.size++
		return
	}

	fetchedNode := m.root.fetchNode(key, m.comparator)
	if fetchedNode.data == nil {
		m.size++
	}
	fetchedNode.data = &value
}

func (m *OrderedMap[K, V]) Erase(key K) {
	if m.root == nil {
		return
	}

	if m.root.next == nil && m.root.prev == nil {
		m.root = nil
		m.size--
		return
	}

	nodeToDelete := m.root.findNode(key, m.comparator)
	if nodeToDelete == nil {
		return
	}

	if nodeToDelete.prev == nil {
		nodeToDelete.next.prev = nil
		m.size--
		return
	}
	if nodeToDelete.next == nil {
		nodeToDelete.prev.next = nil
		m.size--
		return
	}

	nodeToDelete.next.prev = nodeToDelete.prev
	nodeToDelete.prev.next = nodeToDelete.next
	m.size--
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	if m.root == nil {
		return false
	}
	dataNode := m.root.findNode(key, m.comparator)
	return dataNode != nil
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size // need to implement
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	if m.root == nil {
		return
	}
	firstNode := m.root.first()
	if firstNode == nil {
		return
	}
	currentNode := firstNode
	for currentNode != nil {
		value := *currentNode.data
		action(currentNode.key, value)
		currentNode = currentNode.next
	}
}

func (m *OrderedMap[K, V]) PrintTree() {
	firstNode := m.root.first()
	currentNode := firstNode
	for currentNode != nil {
		currentKey := currentNode.key
		var prevKey K
		var nextKey K
		if currentNode.next != nil {
			nextKey = currentNode.next.key
		}
		if currentNode.prev != nil {
			prevKey = currentNode.prev.key
		}
		fmt.Println("current key:", currentKey, "prev:", prevKey, "next:", nextKey)

		currentNode = currentNode.next
	}
}

func TestCircularQueue(t *testing.T) {
	less := func(a, b int) bool {
		return a > b
	}
	data := NewOrderedMap[int, int](less)
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

func TestInsertDuplicateKey(t *testing.T) {
	less := func(a, b int) bool {
		return a > b
	}
	m := NewOrderedMap[int, int](less)
	m.Insert(1, 10)
	m.Insert(1, 20) // duplicate key, should update value, not size
	assert.Equal(t, 1, m.Size())
	var val int
	m.ForEach(func(_, v int) { val = v })
	assert.Equal(t, 20, val)
}

func TestEraseNonExistent(t *testing.T) {
	less := func(a, b int) bool {
		return a > b
	}
	m := NewOrderedMap[int, int](less)
	m.Insert(1, 1)
	m.Insert(2, 2)
	m.Erase(3) // erase non-existent key
	assert.Equal(t, 2, m.Size())
}

func TestEraseFromEmpty(t *testing.T) {
	less := func(a, b int) bool {
		return a > b
	}
	m := NewOrderedMap[int, int](less)
	m.Erase(1) // should not panic
	assert.Zero(t, m.Size())
}

func TestInsertEraseSingle(t *testing.T) {
	less := func(a, b int) bool {
		return a > b
	}
	m := NewOrderedMap[int, int](less)
	m.Insert(42, 100)
	assert.Equal(t, 1, m.Size())
	m.Erase(42)
	assert.Zero(t, m.Size())
	assert.False(t, m.Contains(42))
}

func TestForEachEmpty(t *testing.T) {
	less := func(a, b int) bool {
		return a > b
	}
	m := NewOrderedMap[int, int](less)
	called := false
	m.ForEach(func(_, _ int) { called = true })
	assert.False(t, called)
}

func TestOrderAfterMixedOps(t *testing.T) {
	less := func(a, b int) bool {
		return a > b
	}
	m := NewOrderedMap[int, int](less)
	m.Insert(5, 50)
	m.Insert(1, 10)
	m.Insert(3, 30)
	m.Insert(2, 20)
	m.Insert(4, 40)
	m.Erase(3)
	m.Insert(3, 300)
	m.Erase(1)
	m.Insert(6, 60)
	var keys []int
	m.ForEach(func(k, _ int) { keys = append(keys, k) })
	expected := []int{2, 3, 4, 5, 6}
	assert.True(t, reflect.DeepEqual(expected, keys))
}

func TestNegativeAndLargeKeys(t *testing.T) {
	less := func(a, b int) bool {
		return a > b
	}
	m := NewOrderedMap[int, int](less)
	m.Insert(-10, 1)
	m.Insert(0, 2)
	m.Insert(1000000, 3)
	m.Insert(-100, 4)
	var keys []int
	m.ForEach(func(k, _ int) { keys = append(keys, k) })
	expected := []int{-100, -10, 0, 1000000}
	assert.True(t, reflect.DeepEqual(expected, keys))
}

func TestMinMaxIntKeys(t *testing.T) {
	less := func(a, b int) bool {
		return a > b
	}
	m := NewOrderedMap[int, int](less)
	min := -1 << 63
	max := 1<<63 - 1
	m.Insert(min, 111)
	m.Insert(max, 222)
	m.Insert(0, 333)
	var keys []int
	m.ForEach(func(k, _ int) { keys = append(keys, k) })
	expected := []int{min, 0, max}
	assert.True(t, reflect.DeepEqual(expected, keys))
	assert.True(t, m.Contains(min))
	assert.True(t, m.Contains(max))
	assert.True(t, m.Contains(0))
}

func TestStringKeysAndValues(t *testing.T) {
	less := func(a, b string) bool { return a > b }
	m := NewOrderedMap[string, string](less)
	m.Insert("b", "bee")
	m.Insert("a", "alpha")
	m.Insert("c", "cat")
	var keys []string
	var values []string
	m.ForEach(func(k, v string) {
		keys = append(keys, k)
		values = append(values, v)
	})
	assert.True(t, reflect.DeepEqual([]string{"a", "b", "c"}, keys))
	assert.True(t, reflect.DeepEqual([]string{"alpha", "bee", "cat"}, values))
}

func TestFloat64Keys(t *testing.T) {
	less := func(a, b float64) bool { return a > b }
	m := NewOrderedMap[float64, int](less)
	m.Insert(2.2, 22)
	m.Insert(1.1, 11)
	m.Insert(3.3, 33)
	var keys []float64
	m.ForEach(func(k float64, _ int) { keys = append(keys, k) })
	assert.True(t, reflect.DeepEqual([]float64{1.1, 2.2, 3.3}, keys))
}

func TestRuneKeys(t *testing.T) {
	less := func(a, b rune) bool { return a > b }
	m := NewOrderedMap[rune, string](less)
	m.Insert('b', "bee")
	m.Insert('a', "alpha")
	m.Insert('c', "cat")
	var keys []rune
	m.ForEach(func(k rune, _ string) { keys = append(keys, k) })
	assert.True(t, reflect.DeepEqual([]rune{'a', 'b', 'c'}, keys))
}

type person struct {
	Name string
	Age  int
}

func TestStructValues(t *testing.T) {
	less := func(a, b int) bool { return a > b }
	m := NewOrderedMap[int, person](less)
	m.Insert(2, person{"Bob", 30})
	m.Insert(1, person{"Alice", 25})
	m.Insert(3, person{"Carol", 40})
	var names []string
	m.ForEach(func(_ int, v person) { names = append(names, v.Name) })
	assert.True(t, reflect.DeepEqual([]string{"Alice", "Bob", "Carol"}, names))
}

type keyByField struct {
	ID  int
	Tag string
}

func TestStructKeysByField(t *testing.T) {
	less := func(a, b keyByField) bool { return a.ID > b.ID }
	m := NewOrderedMap[keyByField, string](less)
	m.Insert(keyByField{2, "b"}, "bee")
	m.Insert(keyByField{1, "a"}, "alpha")
	m.Insert(keyByField{3, "c"}, "cat")
	var ids []int
	m.ForEach(func(k keyByField, _ string) { ids = append(ids, k.ID) })
	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, ids))
}
