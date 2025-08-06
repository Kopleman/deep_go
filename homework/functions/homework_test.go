package main

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Map[S, T any](data []S, action func(S) T) []T {
	if data == nil {
		return nil
	}
	if len(data) == 0 {
		return []T{}
	}
	ret := make([]T, len(data))
	for i := 0; i < len(data); i++ {
		ret[i] = action(data[i])
	}
	return ret
}

func Filter[T any](data []T, action func(T) bool) []T {
	if len(data) == 0 {
		return data
	}
	result := make([]T, 0, len(data))
	for i := 0; i < len(data); i++ {
		if action(data[i]) {
			result = append(result, data[i])
		}
	}
	return result
}

func Reduce[T any](data []T, initial T, action func(T, T) T) T {
	if len(data) == 0 {
		return initial
	}
	for i := 0; i < len(data); i++ {
		initial = action(initial, data[i])
	}
	return initial
}

func TestMap(t *testing.T) {
	tests := map[string]struct {
		data   []int
		action func(int) int
		result []int
	}{
		"nil numbers": {
			action: func(number int) int {
				return -number
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(number int) int {
				return -number
			},
			result: []int{},
		},
		"inc numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) int {
				return number + 1
			},
			result: []int{2, 3, 4, 5, 6},
		},
		"double numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) int {
				return number * number
			},
			result: []int{1, 4, 9, 16, 25},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Map(test.data, test.action)
			spew.Dump(result)
			assert.True(t, reflect.DeepEqual(test.result, result))
		})
	}
}

func TestFilter(t *testing.T) {
	tests := map[string]struct {
		data   []int
		action func(int) bool
		result []int
	}{
		"nil numbers": {
			action: func(number int) bool {
				return number == 0
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(number int) bool {
				return number == 1
			},
			result: []int{},
		},
		"even numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) bool {
				return number%2 == 0
			},
			result: []int{2, 4},
		},
		"positive numbers": {
			data: []int{-1, -2, 1, 2},
			action: func(number int) bool {
				return number > 0
			},
			result: []int{1, 2},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Filter(test.data, test.action)
			assert.True(t, reflect.DeepEqual(test.result, result))
		})
	}
}

func TestReduce(t *testing.T) {
	tests := map[string]struct {
		initial int
		data    []int
		action  func(int, int) int
		result  int
	}{
		"nil numbers": {
			action: func(lhs, rhs int) int {
				return 0
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(lhs, rhs int) int {
				return 0
			},
		},
		"sum of numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(lhs, rhs int) int {
				return lhs + rhs
			},
			result: 15,
		},
		"sum of numbers with initial value": {
			initial: 10,
			data:    []int{1, 2, 3, 4, 5},
			action: func(lhs, rhs int) int {
				return lhs + rhs
			},
			result: 25,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Reduce(test.data, test.initial, test.action)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestMap_Extended(t *testing.T) {
	t.Run("map strings to their lengths", func(t *testing.T) {
		data := []string{"go", "lang", "test"}
		result := Map(data, func(s string) int { return len(s) })
		assert.Equal(t, []int{2, 4, 4}, result)
	})

	t.Run("map struct to string", func(t *testing.T) {
		type S struct{ V int }
		data := []S{{1}, {2}, {3}}
		result := Map(data, func(s S) string { return "val" })
		assert.Equal(t, []string{"val", "val", "val"}, result)
	})

	t.Run("map slice of slices to their lengths", func(t *testing.T) {
		data := [][]int{{1, 2}, {}, {3}}
		result := Map(data, func(s []int) int { return len(s) })
		assert.Equal(t, []int{2, 0, 1}, result)
	})

	t.Run("map with identity function", func(t *testing.T) {
		data := []int{1, 2, 3}
		result := Map(data, func(x int) int { return x })
		assert.Equal(t, data, result)
	})

	t.Run("map with constant function", func(t *testing.T) {
		data := []int{1, 2, 3}
		result := Map(data, func(x int) int { return 42 })
		assert.Equal(t, []int{42, 42, 42}, result)
	})

	t.Run("map with nil action panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic with nil action")
			}
		}()
		_ = Map[int, int]([]int{1}, nil)
	})

	t.Run("copy-on-write: changing input after map does not affect result", func(t *testing.T) {
		data := []int{1, 2, 3}
		result := Map(data, func(x int) int { return x + 1 })
		data[0] = 100
		assert.Equal(t, []int{2, 3, 4}, result)
	})
}

func TestFilter_Extended(t *testing.T) {
	t.Run("filter strings by length", func(t *testing.T) {
		data := []string{"go", "lang", "a"}
		result := Filter(data, func(s string) bool { return len(s) > 1 })
		assert.Equal(t, []string{"go", "lang"}, result)
	})

	t.Run("filter struct by field", func(t *testing.T) {
		type S struct{ V int }
		data := []S{{1}, {2}, {3}}
		result := Filter(data, func(s S) bool { return s.V%2 == 1 })
		assert.Equal(t, []S{{1}, {3}}, result)
	})

	t.Run("filter slice of slices by length", func(t *testing.T) {
		data := [][]int{{1, 2}, {}, {3}}
		result := Filter(data, func(s []int) bool { return len(s) > 0 })
		assert.Equal(t, [][]int{{1, 2}, {3}}, result)
	})

	t.Run("filter with always true", func(t *testing.T) {
		data := []int{1, 2, 3}
		result := Filter(data, func(x int) bool { return true })
		assert.Equal(t, data, result)
	})

	t.Run("filter with always false", func(t *testing.T) {
		data := []int{1, 2, 3}
		result := Filter(data, func(x int) bool { return false })
		assert.Equal(t, []int{}, result)
	})

	t.Run("filter with nil action panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic with nil action")
			}
		}()
		_ = Filter([]int{1}, nil)
	})

	t.Run("copy-on-write: changing input after filter does not affect result", func(t *testing.T) {
		data := []int{1, 2, 3}
		result := Filter(data, func(x int) bool { return x > 1 })
		data[1] = 100
		assert.Equal(t, []int{2, 3}, result)
	})
}

func TestReduce_Extended(t *testing.T) {
	t.Run("reduce strings to concatenation", func(t *testing.T) {
		data := []string{"a", "b", "c"}
		result := Reduce(data, "", func(a, b string) string { return a + b })
		assert.Equal(t, "abc", result)
	})

	t.Run("reduce struct to count (via int slice)", func(t *testing.T) {
		type S struct{ V int }
		data := []S{{1}, {2}, {3}}
		intData := Map(data, func(s S) int { return 1 })
		result := Reduce(intData, 0, func(a, b int) int { return a + b })
		assert.Equal(t, 3, result)
	})

	t.Run("reduce slice of slices to total length (via int slice)", func(t *testing.T) {
		data := [][]int{{1, 2}, {}, {3}}
		lengths := Map(data, func(s []int) int { return len(s) })
		result := Reduce(lengths, 0, func(a, b int) int { return a + b })
		assert.Equal(t, 3, result)
	})

	t.Run("reduce with idempotent function", func(t *testing.T) {
		data := []int{1, 2, 3}
		result := Reduce(data, 0, func(a, b int) int { return a })
		assert.Equal(t, 0, result)
	})

	t.Run("reduce with constant function", func(t *testing.T) {
		data := []int{1, 2, 3}
		result := Reduce(data, 0, func(a, b int) int { return 42 })
		assert.Equal(t, 42, result)
	})

	t.Run("reduce with nil action panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic with nil action")
			}
		}()
		_ = Reduce([]int{1}, 0, nil)
	})

	t.Run("copy-on-write: changing input after reduce does not affect result", func(t *testing.T) {
		data := []int{1, 2, 3}
		result := Reduce(data, 0, func(a, b int) int { return a + b })
		data[0] = 100
		assert.Equal(t, 6, result)
	})
}
