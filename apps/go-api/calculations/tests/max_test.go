package calculations

import (
	"api/calculations"
	"reflect"
	"testing"
)

func TestMax(t *testing.T) {
	tests := []struct {
		input    []int
		expected int
	}{
		{[]int{1, 2, 3}, 3},
		{[]int{-1, -2, -3}, -1},
		{[]int{0, 0, 0}, 0},
	}
	for _, tt := range tests {
		result := calculations.Max(tt.input)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("Max(%v) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
