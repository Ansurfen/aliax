// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package aslices

// Map applies a function to each element in a slice in-place.
func MapInPlace[S ~[]E, E any](s S, f func(E) E) {
	for i := range s {
		s[i] = f(s[i])
	}
}

// Map applies a function to each element in a slice and returns a new slice.
func Map[S ~[]E, E any, R any](s S, f func(E) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}
