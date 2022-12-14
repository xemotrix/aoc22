package evilgo

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Pipe pipes two functions so: Pipe(A, B)(c) == B(A(c))
func Pipe[T, K, R any](a func(T) K, b func(K) R) func(T) R {
	return func(input T) R {
		return b(a(input))
	}

}

// The identity function Identity(x) == x
func Identity[T any](i T) T {
	return i
}

// Map maps a function to each elent of an slice and returns a new slice with
// the results
func Map[T, R any](function func(T) R, slice []T) []R {
	result := make([]R, len(slice))
	for i, ele := range slice {

		result[i] = function(ele)
	}
	return result
}

func Apply[T, R any](function func(T) R, val T) R {
	return function(val)
}

// Reduce applies a function n times to convert an slice of T into a single
// value of type T
func Reduce[T any](function func(T, T) T, slice []T) T {
	if len(slice) == 0 {
		panic("Can't reduce if slice is empty")
	}
	res := slice[0]
	for _, ele := range slice[1:] {
		res = function(res, ele)
	}
	return res
}

// Filp swaps the inputs of a dyadic function and returns a new function
func Flip[T, K, R any](function func(T, K) R) func(K, T) R {
	return func(a K, b T) R {
		return function(b, a)
	}

}

func Scan[T any](function func(T, T) T, slice []T) []T {
	if len(slice) <= 1 {
		return slice
	}
	res := make([]T, len(slice))
	res[0] = slice[0]
	for i, ele := range slice[1:] {
		res[i+1] = function(res[i], ele)
	}
	return res
}

// Curry takes a dyadic function and its first argument and returns a monadic
// function that expects the second argument
func Curry[T, K, R any](function func(T, K) R, part1 T) func(K) R {
	return func(part2 K) R {
		return function(part1, part2)
	}
}

// Curry takes a dyadic function and its second argument and returns a monadic
// function that expects the first argument
func FCurry[T, K, R any](function func(T, K) R, part1 K) func(T) R {
	return func(part2 T) R {
		return Flip(function)(part1, part2)
	}
}

// Expect takes a function that returns (R, error) and returns a function
// that returns R and panics if the original function would have retuned an
// error
func Expect[T, R any](function func(T) (R, error)) func(T) R {
	return func(in T) R {
		res, err := function(in)
		if err != nil {
			panic(fmt.Sprintf("Expected no error but got error: %s", err))
		}
		return res
	}
}

// Filter returns every element of the input slice that returns true when
// evaluated with the imput function
func Filter[T any](function func(T) bool, slice []T) (res []T) {
	for _, ele := range slice {
		if function(ele) {
			res = append(res, ele)
		}
	}
	return
}

// Max takes two arguments of a type T that satisfies the type constraint
// Ordered and returns the greater one
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min takes two arguments of a type T that satisfies the type constraint
// Ordered and returns the greater one
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
