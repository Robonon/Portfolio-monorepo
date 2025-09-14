package calculations

import (
	"api/calculations"
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{[]int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{1}, []int{1}},
		{[]int{}, []int{}},
	}

	for _, tt := range tests {
		result := calculations.Reverse(append([]int{}, tt.input...)) // copy input to avoid mutation
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("Reverse(%v) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
