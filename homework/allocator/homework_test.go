package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	if len(memory) == 0 {
		return
	}

	memoryStart := uintptr(unsafe.Pointer(&memory[0]))
	memoryEnd := uintptr(unsafe.Pointer(&memory[len(memory)-1]))

	for _, ptr := range pointers {
		if ptr == nil {
			panic("nil pointer found")
		}
		ptrAddr := uintptr(ptr)
		if ptrAddr < memoryStart || ptrAddr > memoryEnd {
			panic("pointer outside memory bounds")
		}
		bytePtr := (*byte)(ptr)
		if *bytePtr == 0x00 {
			panic("pointer points to free byte")
		}
	}

	writePos := 0
	readPos := 0
	pointerIndex := 0

	for readPos < len(memory) {
		currentPtr := unsafe.Pointer(&memory[readPos])

		if pointerIndex < len(pointers) && currentPtr == pointers[pointerIndex] {
			if writePos != readPos {
				memory[writePos] = memory[readPos]
				pointers[pointerIndex] = unsafe.Pointer(&memory[writePos])
			}
			writePos++
			pointerIndex++
		}
		readPos++
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

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}

func TestDefragmentationEmptyMemory(t *testing.T) {
	memory := []byte{}
	pointers := []unsafe.Pointer{}

	// Не должно паниковать
	Defragment(memory, pointers)
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

	Defragment(memory, pointers)
	assert.True(t, reflect.DeepEqual(expectedMemory, memory))
	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}

func TestDefragmentationNilPointer(t *testing.T) {
	memory := []byte{0xFF, 0x00, 0xFF, 0x00}
	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		nil, // nil pointer
		unsafe.Pointer(&memory[2]),
	}

	assert.Panics(t, func() {
		Defragment(memory, pointers)
	})
}

func TestDefragmentationPointerToFreeByte(t *testing.T) {
	memory := []byte{0xFF, 0x00, 0xFF, 0x00}
	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		unsafe.Pointer(&memory[1]), // указывает на 0x00
		unsafe.Pointer(&memory[2]),
	}

	assert.Panics(t, func() {
		Defragment(memory, pointers)
	})
}

func TestDefragmentationComplexCase(t *testing.T) {
	memory := []byte{
		0xAA, 0x00, 0xBB, 0x00, 0xCC, 0x00, 0xDD, 0x00,
		0x00, 0xEE, 0x00, 0xFF, 0x00, 0x11, 0x00, 0x22,
	}

	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),  // 0xAA
		unsafe.Pointer(&memory[2]),  // 0xBB
		unsafe.Pointer(&memory[4]),  // 0xCC
		unsafe.Pointer(&memory[6]),  // 0xDD
		unsafe.Pointer(&memory[9]),  // 0xEE
		unsafe.Pointer(&memory[11]), // 0xFF
		unsafe.Pointer(&memory[13]), // 0x11
		unsafe.Pointer(&memory[15]), // 0x22
	}

	expectedMemory := []byte{
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	Defragment(memory, pointers)
	assert.True(t, reflect.DeepEqual(expectedMemory, memory))

	// Проверяем, что указатели указывают на правильные позиции
	for i := 0; i < len(pointers); i++ {
		bytePtr := (*byte)(pointers[i])
		assert.Equal(t, expectedMemory[i], *bytePtr)
	}
}

func BenchmarkDefragmentation(b *testing.B) {
	// Создаем большой фрагментированный массив
	size := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Создаем новые массивы для каждого теста
		memory := make([]byte, size)
		pointers := make([]unsafe.Pointer, size/2)

		// Заполняем память чередующимися занятыми и свободными байтами
		for j := 0; j < size; j++ {
			if j%2 == 0 {
				memory[j] = 0xFF // занятый байт (всегда не 0x00)
			} else {
				memory[j] = 0x00 // свободный байт
			}
		}

		// Создаем указатели на занятые байты
		pointerIndex := 0
		for j := 0; j < size && pointerIndex < len(pointers); j++ {
			if j%2 == 0 {
				pointers[pointerIndex] = unsafe.Pointer(&memory[j])
				pointerIndex++
			}
		}

		// Обрезаем массив указателей до фактического размера
		testPointers := pointers[:pointerIndex]

		Defragment(memory, testPointers)
	}
}
