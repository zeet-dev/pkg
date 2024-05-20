package utils

import "github.com/samber/lo"

func Ptr[T any](in T) *T {
	return &in
}

// Count returns the count of true values
func Count(v ...bool) int {
	return lo.Count(v, true)
}

// MapPtr runs f on the pointer value if it is not nil, or returns nil
func MapPtr[I, O any](in *I, f func(I) O) *O {
	if in != nil {
		return Ptr(f(*in))
	}
	return nil
}

func PtrSlice[T any](in []T) []*T {
	out := make([]*T, len(in))
	for i := range in {
		out[i] = &in[i]
	}
	return out
}

func ValueSlice[T any](in []*T) []T {
	out := make([]T, len(in))
	for i := range in {
		out[i] = *in[i]
	}
	return out
}
