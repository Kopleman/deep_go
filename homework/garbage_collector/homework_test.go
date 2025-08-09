package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Trace(stacks [][]uintptr) []uintptr {
	visited := make(map[*uintptr]bool)
	var result []uintptr

	for _, s := range stacks {
		for _, ptr := range s {
			memPtr := (*uintptr)(unsafe.Pointer(ptr))
			if ptr != 0 && !visited[memPtr] {
				visited[memPtr] = true
				result = append(result, ptr)

				processObject(ptr, visited, &result)
			}
		}
	}

	return result
}

func processObject(ptr uintptr, visited map[*uintptr]bool, result *[]uintptr) {
	memPtr := (*uintptr)(unsafe.Pointer(ptr))
	value := *memPtr
	valueAsMemPointer := (*uintptr)(unsafe.Pointer(value))
	if value != 0 && !visited[valueAsMemPointer] {
		visited[valueAsMemPointer] = true
		*result = append(*result, value)

		processObject(value, visited, result)
	}
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
	}

	pointers := Trace(stacks)
	/**
	[0]heapPointer1->[1]heapObjects[1]
	и только потом мы дойдем до
	[2]heapObjects[0]
	*/
	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
	}

	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}

func TestTraceEmptyStacks(t *testing.T) {
	stacks := [][]uintptr{}
	pointers := Trace(stacks)
	assert.Empty(t, pointers)

	stacks = [][]uintptr{{}, {}, {}}
	pointers = Trace(stacks)
	assert.Empty(t, pointers)
}

func TestTraceOnlyZeros(t *testing.T) {
	stacks := [][]uintptr{
		{0x00, 0x00, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00},
	}
	pointers := Trace(stacks)
	assert.Empty(t, pointers)
}

func TestTraceCircularReferences(t *testing.T) {
	var obj1, obj2 int
	var ptr1 *int = &obj1
	var ptr2 *int = &obj2

	// Создаем циклическую ссылку через unsafe
	*(*uintptr)(unsafe.Pointer(&obj1)) = uintptr(unsafe.Pointer(&obj2))
	*(*uintptr)(unsafe.Pointer(&obj2)) = uintptr(unsafe.Pointer(&obj1))

	stacks := [][]uintptr{
		{uintptr(unsafe.Pointer(&ptr1))},
		{uintptr(unsafe.Pointer(&ptr2))},
	}

	pointers := Trace(stacks)

	expectedCount := 4
	assert.Len(t, pointers, expectedCount)

	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&ptr1)),
		uintptr(unsafe.Pointer(&obj1)),
		uintptr(unsafe.Pointer(&ptr2)),
		uintptr(unsafe.Pointer(&obj2)),
	}

	for _, expected := range expectedPointers {
		assert.Contains(t, pointers, expected)
	}
}

func TestTraceSelfReference(t *testing.T) {
	var obj int
	var ptr *int = &obj

	*(*uintptr)(unsafe.Pointer(&obj)) = uintptr(unsafe.Pointer(&obj))

	stacks := [][]uintptr{
		{uintptr(unsafe.Pointer(&ptr))},
	}

	pointers := Trace(stacks)

	expectedCount := 2
	assert.Len(t, pointers, expectedCount)

	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&ptr)),
		uintptr(unsafe.Pointer(&obj)),
	}

	for _, expected := range expectedPointers {
		assert.Contains(t, pointers, expected)
	}
}

func TestTraceDeepChain(t *testing.T) {
	var objects [5]int
	var ptrArray [5]*int

	for i := 0; i < 5; i++ {
		ptrArray[i] = &objects[i]
		if i < 4 {
			*(*uintptr)(unsafe.Pointer(&objects[i])) = uintptr(unsafe.Pointer(&ptrArray[i+1]))
		}
	}

	stacks := [][]uintptr{
		{uintptr(unsafe.Pointer(&ptrArray[0]))},
	}

	result := Trace(stacks)

	expectedCount := 10
	assert.Len(t, result, expectedCount)

	for i := 0; i < 5; i++ {
		assert.Contains(t, result, uintptr(unsafe.Pointer(&ptrArray[i])))
		assert.Contains(t, result, uintptr(unsafe.Pointer(&objects[i])))
	}
}

func TestTraceDuplicatePointers(t *testing.T) {
	var obj int
	var ptr *int = &obj

	stacks := [][]uintptr{
		{uintptr(unsafe.Pointer(&ptr)), uintptr(unsafe.Pointer(&ptr))},
		{uintptr(unsafe.Pointer(&ptr))},
		{uintptr(unsafe.Pointer(&obj)), uintptr(unsafe.Pointer(&ptr))},
	}

	pointers := Trace(stacks)

	expectedCount := 2
	assert.Len(t, pointers, expectedCount)

	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&ptr)),
		uintptr(unsafe.Pointer(&obj)),
	}

	for _, expected := range expectedPointers {
		assert.Contains(t, pointers, expected)
	}
}

func TestTraceNilPointers(t *testing.T) {
	var nilPtr *int = nil
	var obj int
	var validPtr *int = &obj

	stacks := [][]uintptr{
		{uintptr(unsafe.Pointer(&nilPtr)), uintptr(unsafe.Pointer(&validPtr))},
		{0x00, uintptr(unsafe.Pointer(&obj))},
	}

	pointers := Trace(stacks)

	expectedCount := 3
	assert.Len(t, pointers, expectedCount)

	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&nilPtr)),
		uintptr(unsafe.Pointer(&validPtr)),
		uintptr(unsafe.Pointer(&obj)),
	}

	for _, expected := range expectedPointers {
		assert.Contains(t, pointers, expected)
	}
}

func TestTraceMixedData(t *testing.T) {
	var objects [3]int
	var ptrArray [3]*int

	for i := 0; i < 3; i++ {
		ptrArray[i] = &objects[i]
	}

	stacks := [][]uintptr{
		{
			uintptr(unsafe.Pointer(&ptrArray[0])), 0x00, 0x00,
			uintptr(unsafe.Pointer(&objects[0])), 0x00, 0x00,
		},
		{
			0x00, uintptr(unsafe.Pointer(&ptrArray[1])), 0x00,
			uintptr(unsafe.Pointer(&objects[1])), 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&ptrArray[2])), 0x00, 0x00,
			0x00, uintptr(unsafe.Pointer(&objects[2])), 0x00,
		},
	}

	result := Trace(stacks)

	expectedCount := 6
	assert.Len(t, result, expectedCount)

	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&ptrArray[0])),
		uintptr(unsafe.Pointer(&objects[0])),
		uintptr(unsafe.Pointer(&ptrArray[1])),
		uintptr(unsafe.Pointer(&objects[1])),
		uintptr(unsafe.Pointer(&ptrArray[2])),
		uintptr(unsafe.Pointer(&objects[2])),
	}

	for _, expected := range expectedPointers {
		assert.Contains(t, result, expected)
	}
}
