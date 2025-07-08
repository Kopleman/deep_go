package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type COWBuffer struct {
	data []byte
	refs *int
}

func NewCOWBuffer(data []byte) COWBuffer {
	return COWBuffer{
		data: data,
		refs: new(int),
	}
}

func (b *COWBuffer) Clone() COWBuffer {
	*b.refs++
	return COWBuffer{
		data: b.data,
		refs: b.refs,
	}
}

func (b *COWBuffer) Close() {
	*b.refs--
	if *b.refs <= 0 {
		b.data = nil
	}
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if index < 0 || index >= len(b.data) {
		return false
	}

	/*
		1. we reduce ref number for current siblings
		2. allocate new array
		3. create new pointer for refs (destroy connection between siblings)
	*/
	if *b.refs > 0 {
		*b.refs--
		b.data = append([]byte(nil), b.data...)
		b.refs = new(int)
	}

	b.data[index] = value
	return true
}

func (b *COWBuffer) String() string {
	if len(b.data) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b.data), len(b.data))
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
}

func TestCOWBuffer_EdgeCases(t *testing.T) {
	// Empty buffer
	empty := NewCOWBuffer([]byte{})
	defer empty.Close()
	assert.Equal(t, "", empty.String())
	assert.False(t, empty.Update(0, 'x'))
	assert.False(t, empty.Update(-1, 'x'))
	assert.False(t, empty.Update(1, 'x'))

	// Update first and last element
	buf := NewCOWBuffer([]byte{'x', 'y', 'z'})
	defer buf.Close()
	assert.True(t, buf.Update(0, 'a'))
	assert.True(t, buf.Update(2, 'b'))
	assert.Equal(t, []byte{'a', 'y', 'b'}, buf.data)
}

func TestCOWBuffer_CopyOnWriteIsolation(t *testing.T) {
	buf := NewCOWBuffer([]byte{'1', '2', '3'})
	defer buf.Close()
	clone := buf.Clone()
	defer clone.Close()

	// Change original — clone should not change
	buf.Update(1, 'x')
	assert.Equal(t, []byte{'1', 'x', '3'}, buf.data)
	assert.Equal(t, []byte{'1', '2', '3'}, clone.data)

	// Change clone — original should not change
	clone.Update(2, 'y')
	assert.Equal(t, []byte{'1', 'x', '3'}, buf.data)
	assert.Equal(t, []byte{'1', '2', 'y'}, clone.data)
}

func TestCOWBuffer_RefsAndClose(t *testing.T) {
	buf := NewCOWBuffer([]byte{'a', 'b'})
	clone1 := buf.Clone()
	clone2 := buf.Clone()

	assert.Equal(t, 2, *buf.refs) // 2 clones
	clone1.Close()
	assert.Equal(t, 1, *buf.refs)
	clone2.Close()
	assert.Equal(t, 0, *buf.refs)
	buf.Close()
	assert.Nil(t, buf.data)

	// Repeated Close should not panic
	buf.Close()
}

func TestCOWBuffer_StringUnicode(t *testing.T) {
	// UTF-8 symbols
	data := []byte("Привет")
	buf := NewCOWBuffer(data)
	defer buf.Close()
	assert.Equal(t, "Привет", buf.String())
}

func TestCOWBuffer_MultipleClonesAndUpdates(t *testing.T) {
	buf := NewCOWBuffer([]byte{'a', 'b', 'c'})
	defer buf.Close()
	clones := make([]COWBuffer, 5)
	for i := range clones {
		clones[i] = buf.Clone()
		defer clones[i].Close()
	}
	// Update each clone separately
	for i := range clones {
		clones[i].Update(i%3, byte('x'+i))
	}
	// Original should not change
	assert.Equal(t, []byte{'a', 'b', 'c'}, buf.data)
	// Each clone should be unique
	for i := range clones {
		for j := range clones {
			if i != j {
				assert.NotEqual(t, clones[i].data, clones[j].data)
			}
		}
	}
}
