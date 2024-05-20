package utils

import (
	"reflect"

	"github.com/samber/lo"
	"go.uber.org/multierr"
)

func ReverseSlice(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func SliceContains[S ~[]T, T comparable](items S, el T) bool {
	for _, val := range items {
		if val == el {
			return true
		}
	}
	return false
}

func SliceOfPointersToValue[S []*E, E any](s S) []E {
	var out []E
	for _, el := range s {
		out = append(out, *el)
	}
	return out
}

func SliceOfValuesToPointer[S []E, E any](s S) []*E {
	var out []*E
	for _, el := range s {
		currEl := el
		out = append(out, &currEl)
	}
	return out
}

func SliceFindUniques[T comparable](collection []T) []T {
	isDupl := make(map[T]bool, len(collection))

	for _, item := range collection {
		duplicated, ok := isDupl[item]
		if !ok {
			isDupl[item] = false
		} else if !duplicated {
			isDupl[item] = true
		}
	}

	result := make([]T, 0, len(collection)-len(isDupl))

	for _, item := range collection {
		if duplicated := isDupl[item]; !duplicated {
			result = append(result, item)
		}
	}

	return result
}

// Chain is an alias for MultiAppend
func Chain[T any](slices ...[]T) []T {
	return MultiAppend(slices...)
}

func MultiAppend[T any](slices ...[]T) []T {
	outCapacity := 0
	for _, s := range slices {
		outCapacity += len(s)
	}

	out := make([]T, 0, outCapacity)
	for _, s := range slices {
		out = append(out, s...)
	}

	return out
}

// FilterSliceToMap iterates over a slice, turning T items into (K, V) associations, with possible filtering
// iteratee returns (bool, K, V): when the bool is false, K=>V is not included in the return value
// adapted from samber/lo
func FilterSliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (bool, K, V)) map[K]V {
	result := make(map[K]V)

	for _, v := range collection {
		ok, key, value := transform(v)
		if ok {
			result[key] = value
		}
	}

	return result
}

// TryFilter filters a collection, but can abort when an error is returned by the iteratee
// when an error is returned, the collection return value will be empty
func TryFilter[T any, R any](collection []T, iteratee func(item T, index int) (bool, R, error)) ([]R, error) {
	result := make([]R, 0, len(collection))

	for index, item := range collection {
		if ok, r, err := iteratee(item, index); err != nil {
			return nil, err
		} else if ok {
			result = append(result, r)
		}
	}

	return result, nil
}

// FindOrError returns the first T in collection which matches the predicate, if none match, err is returned
func FindOrError[T any](collection []T, err error, predicate func(item T) bool) (zero T, errNil error) {
	if item, ok := lo.Find(collection, predicate); ok {
		return item, nil
	} else {
		return zero, err
	}
}

// TryMap maps collection through a fallible mapper, returning the mapped collection and any error
// when an error is returned, the failed inputs will have a corresponding zero value in the returned collection
func TryMap[T any, R any](collection []T, mapper func(item T) (R, error)) ([]R, error) {
	var errs []error
	mapped := lo.Map(collection, func(item T, index int) (zero R) {
		if result, err := mapper(item); err != nil {
			errs = append(errs, err)
			return zero
		} else {
			return result
		}
	})
	return mapped, multierr.Combine(errs...)
}
