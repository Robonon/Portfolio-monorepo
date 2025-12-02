package calculations

import (
	"api/calculations"
	"reflect"
	"testing"
)

func TestCountUnique(t *testing.T) {
	tests := []struct {
		input    []int
		expected int
	}{
		{[]int{1, 2, 3, 1, 2, 3}, 3},
		{[]int{1, 1, 1, 1, 1}, 1},
		{[]int{1, 2, 3, 4, 5}, 5},
		{[]int{}, 0},
	}
	for _, tt := range tests {
		result := calculations.CountUnique(tt.input)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("CountUnique(%v) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
