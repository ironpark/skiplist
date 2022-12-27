// Copyright 2011 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package skiplist

// Comparable defines a comparable func.
type Comparable[K any] func(lhs, rhs K) int

// Reverse creates a reversed comparable.
func Reverse[K any](comparable Comparable[K]) Comparable[K] {
	return func(lhs, rhs K) int {
		return -comparable(lhs, rhs)
	}
}
