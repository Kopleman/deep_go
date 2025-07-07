package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func ToLittleEndian[T uint16 | uint32 | uint64](n T) T {
	size := unsafe.Sizeof(n)

	for i, j := 0, int(size)-1; i < j; i, j = i+1, j-1 {
		left := *(*uint8)(unsafe.Add(unsafe.Pointer(&n), i))
		right := *(*uint8)(unsafe.Add(unsafe.Pointer(&n), j))

		*(*uint8)(unsafe.Add(unsafe.Pointer(&n), i)) = right
		*(*uint8)(unsafe.Add(unsafe.Pointer(&n), j)) = left
	}

	return n
}

func TestConversion(t *testing.T) {
	tests32 := map[string]struct {
		number uint32
		result uint32
	}{
		"uint32 test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"uint32 test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"uint32 test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"uint32 test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"uint32 test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
		"uint32 test case #6": {
			number: 0x01234567, // [01][23][45][67]
			result: 0x67452301, // [67][45][23][01]
		},
		"uint32 test case #7": {
			number: 0x01010202,
			result: 0x02020101,
		},
		"uint32 test case #8": {
			number: 0x01,
			result: 0x1000000,
		},
		"uint32 test case #9": {
			number: 0x1,
			result: 0x1000000,
		},
		"uint32 test case #10": {
			number: 0x123,
			result: 0x23010000,
		},
	}

	for name, test := range tests32 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}

	tests16 := map[string]struct {
		number uint16
		result uint16
	}{
		"uint16 test case #1": {
			number: 0x0000,
			result: 0x0000,
		},
		"uint16 test case #2": {
			number: 0xFFFF,
			result: 0xFFFF,
		},
		"uint16 test case #3": {
			number: 0x00FF,
			result: 0xFF00,
		},
		"uint16 test case #4": {
			number: 0x0102,
			result: 0x0201,
		},
		"uint16 test case #5": {
			number: 0x0123,
			result: 0x2301,
		},
		"uint16 test case #6": {
			number: 0x1,
			result: 0x100,
		},
		"uint16 test case #7": {
			number: 0x123,
			result: 0x2301,
		},
		"uint16 test case #8": {
			number: 0x12,
			result: 0x1200,
		},
	}

	for name, test := range tests16 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}

	tests64 := map[string]struct {
		number uint64
		result uint64
	}{
		"uint64 test case #1": {
			number: 0x0000000000000000,
			result: 0x0000000000000000,
		},
		"uint64 test case #2": {
			number: 0xFFFFFFFFFFFFFFFF,
			result: 0xFFFFFFFFFFFFFFFF,
		},
		"uint64 test case #3": {
			number: 0x00FF00FF00FF00FF,
			result: 0xFF00FF00FF00FF00,
		},
		"uint64 test case #4": {
			number: 0x0000FFFF0000FFFF,
			result: 0xFFFF0000FFFF0000,
		},
		"uint64 test case #5": {
			number: 0x0102030405060708,
			result: 0x0807060504030201,
		},
		"uint64 test case #6": {
			number: 0x0123456789101112, // [01][23][45][67][89][10][11][12]
			result: 0x1211108967452301, // [12][11][10][89][67][45][23][01]
		},
		"uint64 test case #7": {
			number: 0x0101020201010202,
			result: 0x0202010102020101,
		},
		"uint32 test case #8": {
			number: 0x01,
			result: 0x100000000000000,
		},
		"uint32 test case #9": {
			number: 0x1,
			result: 0x100000000000000,
		},
		"uint32 test case #10": {
			number: 0x123,
			result: 0x2301000000000000,
		},
	}

	for name, test := range tests64 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
