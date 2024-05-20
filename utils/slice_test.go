package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeet-dev/pkg/utils"
)

func TestMultiAppend(t *testing.T) {
	type testCase[T any] struct {
		name  string
		input [][]T
		want  []T
	}
	tests := []testCase[string]{
		{
			name:  "given nil input, expect empty slice",
			input: nil,
			want:  []string{},
		},
		{
			name: "given empty slice, expect new empty slice",
			input: [][]string{
				{},
			},
			want: []string{},
		},
		{
			name: "given one slice, with one element, expect new slice with same element",
			input: [][]string{
				{"foo"},
			},
			want: []string{"foo"},
		},
		{
			name: "given one slice, with many elements, expect new slice with same elements",
			input: [][]string{
				{"foo", "bar"},
			},
			want: []string{"foo", "bar"},
		},
		{
			name: "given some slices, expect new slice with same elements",
			input: [][]string{
				{"1"},
				{"2"},
			},
			want: []string{"1", "2"},
		},
		{
			name: "given some slices, each with different elements, expect new slice with all elements",
			input: [][]string{
				{"1"},
				{"2", "3", "4"},
			},
			want: []string{"1", "2", "3", "4"},
		},
		{
			name: "given many differing slices, expect new slice with all elements",
			input: [][]string{
				{"1"},
				{"2", "3", "4"},
				{},
				{"5"},
				{"6", "7", "8", "9"},
			},
			want: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.MultiAppend(tt.input...)
			assert.Equal(t, tt.want, got)
			for _, s := range tt.input {
				assert.NotSame(t, s, got)
			}
		})
	}
}
