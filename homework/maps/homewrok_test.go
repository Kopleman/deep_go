package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go
type node struct {
	data int
	key  int
	next *node
	prev *node
}

func (n *node) insertOrUpdateNode(key, val int) (*node, bool) {
	if n.key == key {
		n.data = val
		return n, false
	}
	if n.key > key {
		if n.prev == nil {
			n.prev = &node{key: key, data: val, next: n}
			return n.prev, true
		}

		if n.prev.key < key {
			newNode := &node{key: key, data: val, prev: n.prev, next: n}
			n.prev.next = newNode
			n.prev = newNode
			return n.prev, true
		}

		return n.prev.insertOrUpdateNode(key, val)
	}

	if n.next == nil {
		n.next = &node{key: key, prev: n, data: val}
		return n.next, true
	}

	if n.next.key > key {
		newNode := &node{key: key, prev: n, next: n.next, data: val}
		n.next.prev = newNode
		n.next = newNode
		return n.next, true
	}

	return n.next.insertOrUpdateNode(key, val)
}

func (n *node) findNode(key int) *node {
	if n.key == key {
		return n
	}
	if n.key > key {
		if n.prev == nil {
			return nil
		}

		if n.prev.key < key {
			return nil
		}

		return n.prev.findNode(key)
	}

	if n.next == nil {
		return nil
	}

	if n.next.key > key {
		return nil
	}

	return n.next.findNode(key)
}

func (n *node) first() *node {
	if n.prev == nil {
		return n
	}

	return n.prev.first()
}

func (n *node) last() *node {
	if n.next == nil {
		return n
	}

	return n.next.last()
}

type OrderedMap struct {
	size int
	root *node
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		size: 0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.root = &node{
			data: value,
			key:  key,
		}
		m.size++
		return
	}

	_, isNew := m.root.insertOrUpdateNode(key, value)
	if isNew {
		m.size++
	}
}

func (m *OrderedMap) Erase(key int) {
	if m.root == nil {
		return
	}

	if m.root.next == nil && m.root.prev == nil {
		m.root = nil
		m.size--
		return
	}

	nodeToDelete := m.root.findNode(key)
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

func (m *OrderedMap) Contains(key int) bool {
	if m.root == nil {
		return false
	}
	dataNode := m.root.findNode(key)
	return dataNode != nil
}

func (m *OrderedMap) Size() int {
	return m.size // need to implement
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	if m.root == nil {
		return
	}
	firstNode := m.root.first()
	if firstNode == nil {
		return
	}
	currentNode := firstNode
	for currentNode != nil {
		action(currentNode.key, currentNode.data)
		currentNode = currentNode.next
	}
}

func (m *OrderedMap) PrintTree() {
	firstNode := m.root.first()
	currentNode := firstNode
	for currentNode != nil {
		currentKey := currentNode.key
		prevKey := 0
		nextKey := 0
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
	data := NewOrderedMap()
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
	m := NewOrderedMap()
	m.Insert(1, 10)
	m.Insert(1, 20) // duplicate key, should update value, not size
	assert.Equal(t, 1, m.Size())
	var val int
	m.ForEach(func(_, v int) { val = v })
	assert.Equal(t, 20, val)
}

func TestEraseNonExistent(t *testing.T) {
	m := NewOrderedMap()
	m.Insert(1, 1)
	m.Insert(2, 2)
	m.Erase(3) // erase non-existent key
	assert.Equal(t, 2, m.Size())
}

func TestEraseFromEmpty(t *testing.T) {
	m := NewOrderedMap()
	m.Erase(1) // should not panic
	assert.Zero(t, m.Size())
}

func TestInsertEraseSingle(t *testing.T) {
	m := NewOrderedMap()
	m.Insert(42, 100)
	assert.Equal(t, 1, m.Size())
	m.Erase(42)
	assert.Zero(t, m.Size())
	assert.False(t, m.Contains(42))
}

func TestForEachEmpty(t *testing.T) {
	m := NewOrderedMap()
	called := false
	m.ForEach(func(_, _ int) { called = true })
	assert.False(t, called)
}

func TestOrderAfterMixedOps(t *testing.T) {
	m := NewOrderedMap()
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
	m := NewOrderedMap()
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
	m := NewOrderedMap()
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
