package calculations

import (
	"api/calculations"
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		input    []int
		expected int
	}{{[]int{1, 2, 3}, 6},
		{[]int{-1, -2, -3}, -6},
		{[]int{0, 0, 0}, 0},
		{[]int{}, 0},
	}
	for _, tt := range tests {
		result := calculations.Sum(tt.input)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("Sum(%v) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
