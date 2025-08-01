package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer, blockByteSize int) {
	if len(memory) == 0 {
		return
	}

	if blockByteSize != 1 && blockByteSize != 2 && blockByteSize != 4 && blockByteSize != 8 {
		panic("invalid block byte size - must be 1, 2, 4, or 8")
	}

	if blockByteSize > len(memory) {
		panic("block size greater than memory size")
	}

	if len(memory)%blockByteSize != 0 {
		panic("memory size must be divisible by block size")
	}

	memoryStart := uintptr(unsafe.Pointer(&memory[0]))
	memoryEnd := uintptr(unsafe.Pointer(&memory[len(memory)-blockByteSize]))

	for _, ptr := range pointers {
		if ptr == nil {
			panic("nil pointer found")
		}
		ptrAddr := uintptr(ptr)
		if ptrAddr < memoryStart || ptrAddr > memoryEnd {
			panic("pointer outside memory bounds")
		}

		if (ptrAddr-memoryStart)%uintptr(blockByteSize) != 0 {
			panic("pointer not aligned to block size")
		}

		bytePtr := (*byte)(ptr)
		if *bytePtr == 0x00 {
			panic("pointer points to free block")
		}
	}

	writePos := 0
	readPos := 0
	pointerIndex := 0

	for readPos < len(memory) {
		currentPtr := unsafe.Pointer(&memory[readPos])
		if pointerIndex < len(pointers) && currentPtr == pointers[pointerIndex] {
			if writePos != readPos {
				for i := 0; i < blockByteSize; i++ {
					memory[writePos+i] = memory[readPos+i]
				}
				pointers[pointerIndex] = unsafe.Pointer(&memory[writePos])
			}
			writePos += blockByteSize
			pointerIndex++
		}
		readPos += blockByteSize
	}

	for writePos < len(memory) {
		memory[writePos] = 0x00
		writePos++
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers, 1)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}

func TestDefragmentationEmptyMemory(t *testing.T) {
	memory := []byte{}
	pointers := []unsafe.Pointer{}

	Defragment(memory, pointers, 1)
	assert.Equal(t, 0, len(memory))
}

func TestDefragmentationAlreadyDefragmented(t *testing.T) {
	memory := []byte{0xFF, 0xFF, 0x00, 0x00}
	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		unsafe.Pointer(&memory[1]),
	}

	expectedMemory := []byte{0xFF, 0xFF, 0x00, 0x00}
	expectedPointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		unsafe.Pointer(&memory[1]),
	}

	Defragment(memory, pointers, 1)
	assert.True(t, reflect.DeepEqual(expectedMemory, memory))
	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}

func TestDefragmentationNilPointer(t *testing.T) {
	memory := []byte{0xFF, 0x00, 0xFF, 0x00}
	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		nil,
		unsafe.Pointer(&memory[2]),
	}

	assert.Panics(t, func() {
		Defragment(memory, pointers, 1)
	})
}

func TestDefragmentationPointerToFreeByte(t *testing.T) {
	memory := []byte{0xFF, 0x00, 0xFF, 0x00}
	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		unsafe.Pointer(&memory[1]),
		unsafe.Pointer(&memory[2]),
	}

	assert.Panics(t, func() {
		Defragment(memory, pointers, 1)
	})
}

func TestDefragmentationComplexCase(t *testing.T) {
	memory := []byte{
		0xAA, 0x00, 0xBB, 0x00, 0xCC, 0x00, 0xDD, 0x00,
		0x00, 0xEE, 0x00, 0xFF, 0x00, 0x11, 0x00, 0x22,
	}

	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		unsafe.Pointer(&memory[2]),
		unsafe.Pointer(&memory[4]),
		unsafe.Pointer(&memory[6]),
		unsafe.Pointer(&memory[9]),
		unsafe.Pointer(&memory[11]),
		unsafe.Pointer(&memory[13]),
		unsafe.Pointer(&memory[15]),
	}

	expectedMemory := []byte{
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	Defragment(memory, pointers, 1)
	assert.True(t, reflect.DeepEqual(expectedMemory, memory))

	for i := 0; i < len(pointers); i++ {
		bytePtr := (*byte)(pointers[i])
		assert.Equal(t, expectedMemory[i], *bytePtr)
	}
}

func BenchmarkDefragmentation(b *testing.B) {
	size := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memory := make([]byte, size)
		pointers := make([]unsafe.Pointer, size/2)

		for j := 0; j < size; j++ {
			if j%2 == 0 {
				memory[j] = 0xFF
			} else {
				memory[j] = 0x00
			}
		}

		pointerIndex := 0
		for j := 0; j < size && pointerIndex < len(pointers); j++ {
			if j%2 == 0 {
				pointers[pointerIndex] = unsafe.Pointer(&memory[j])
				pointerIndex++
			}
		}

		testPointers := pointers[:pointerIndex]

		Defragment(memory, testPointers, 1)
	}
}

func TestDefragmentation2ByteBlocks(t *testing.T) {
	memory := []byte{
		0xFF, 0xAA, 0x00, 0x00, 0xBB, 0xCC, 0x00, 0x00,
		0x00, 0x00, 0xDD, 0xEE, 0x00, 0x00, 0xFF, 0x11,
	}

	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		unsafe.Pointer(&memory[4]),
		unsafe.Pointer(&memory[10]),
		unsafe.Pointer(&memory[14]),
	}

	expectedMemory := []byte{
		0xFF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	Defragment(memory, pointers, 2)
	assert.True(t, reflect.DeepEqual(expectedMemory, memory))

	for i := 0; i < len(pointers); i++ {
		bytePtr := (*byte)(pointers[i])
		assert.Equal(t, expectedMemory[i*2], *bytePtr)
	}
}

func TestDefragmentation4ByteBlocks(t *testing.T) {
	memory := []byte{
		0xFF, 0xAA, 0xBB, 0xCC, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0xDD, 0xEE, 0xFF, 0x11,
	}

	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		unsafe.Pointer(&memory[12]),
	}

	expectedMemory := []byte{
		0xFF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	Defragment(memory, pointers, 4)
	assert.True(t, reflect.DeepEqual(expectedMemory, memory))

	for i := 0; i < len(pointers); i++ {
		bytePtr := (*byte)(pointers[i])
		assert.Equal(t, expectedMemory[i*4], *bytePtr)
	}
}

func TestDefragmentation8ByteBlocks(t *testing.T) {
	memory := []byte{
		0xFF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
	}

	expectedMemory := []byte{
		0xFF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	Defragment(memory, pointers, 8)
	assert.True(t, reflect.DeepEqual(expectedMemory, memory))

	bytePtr := (*byte)(pointers[0])
	assert.Equal(t, expectedMemory[0], *bytePtr)
}

func TestDefragmentationInvalidBlockSize(t *testing.T) {
	memory := []byte{0xFF, 0x00, 0xFF, 0x00}
	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		unsafe.Pointer(&memory[2]),
	}

	assert.Panics(t, func() {
		Defragment(memory, pointers, 3)
	})

	assert.Panics(t, func() {
		Defragment(memory, pointers, 16)
	})
}

func TestDefragmentationUnalignedPointer(t *testing.T) {
	memory := []byte{0xFF, 0x00, 0xFF, 0x00}
	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		unsafe.Pointer(&memory[1]),
	}

	assert.Panics(t, func() {
		Defragment(memory, pointers, 2)
	})
}

func TestDefragmentationMemoryNotDivisible(t *testing.T) {
	memory := []byte{0xFF, 0x00, 0xFF}
	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
	}

	assert.Panics(t, func() {
		Defragment(memory, pointers, 2)
	})
}

func BenchmarkDefragmentation2ByteBlocks(b *testing.B) {
	size := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memory := make([]byte, size)
		pointers := make([]unsafe.Pointer, size/4)

		for j := 0; j < size; j += 2 {
			if j%4 == 0 {
				memory[j] = 0xFF
				memory[j+1] = 0xAA
			} else {
				memory[j] = 0x00
				memory[j+1] = 0x00
			}
		}

		pointerIndex := 0
		for j := 0; j < size && pointerIndex < len(pointers); j += 2 {
			if j%4 == 0 {
				pointers[pointerIndex] = unsafe.Pointer(&memory[j])
				pointerIndex++
			}
		}

		testPointers := pointers[:pointerIndex]
		Defragment(memory, testPointers, 2)
	}
}

func BenchmarkDefragmentation4ByteBlocks(b *testing.B) {
	size := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memory := make([]byte, size)
		pointers := make([]unsafe.Pointer, size/8)

		for j := 0; j < size; j += 4 {
			if j%8 == 0 {
				memory[j] = 0xFF
				memory[j+1] = 0xAA
				memory[j+2] = 0xBB
				memory[j+3] = 0xCC
			} else {
				memory[j] = 0x00
				memory[j+1] = 0x00
				memory[j+2] = 0x00
				memory[j+3] = 0x00
			}
		}

		pointerIndex := 0
		for j := 0; j < size && pointerIndex < len(pointers); j += 4 {
			if j%8 == 0 {
				pointers[pointerIndex] = unsafe.Pointer(&memory[j])
				pointerIndex++
			}
		}

		testPointers := pointers[:pointerIndex]
		Defragment(memory, testPointers, 4)
	}
}

func BenchmarkDefragmentation8ByteBlocks(b *testing.B) {
	size := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memory := make([]byte, size)
		pointers := make([]unsafe.Pointer, size/16)

		for j := 0; j < size; j += 8 {
			if j%16 == 0 {
				for k := 0; k < 8; k++ {
					memory[j+k] = byte(0xFF - k)
				}
			} else {
				for k := 0; k < 8; k++ {
					memory[j+k] = 0x00
				}
			}
		}

		pointerIndex := 0
		for j := 0; j < size && pointerIndex < len(pointers); j += 8 {
			if j%16 == 0 {
				pointers[pointerIndex] = unsafe.Pointer(&memory[j])
				pointerIndex++
			}
		}

		testPointers := pointers[:pointerIndex]
		Defragment(memory, testPointers, 8)
	}
}
