// Copyright 2011 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package skiplist

// Element is an element node of a skip list.
type Element[K, V any] struct {
	*elementHeader[K, V]
	Value V
	key   K
	prev  *Element[K, V] // Points to previous adjacent elem.
	list  SkipList[K, V] // The list contains this elem.
	skip  int
}

// elementHeader is the header of an element or a skip list.
// It must be the first anonymous field in a type to make Element() work correctly.
type elementHeader[K, V any] struct {
	next []*Element[K, V] // Next element at all next.
}

// Next returns next adjacent elem.
func (elem *Element[K, V]) Next() *Element[K, V] {
	if len(elem.next) == 0 {
		return nil
	}
	return elem.next[0]
}

// Prev returns previous adjacent elem.
func (elem *Element[K, V]) Prev() *Element[K, V] {
	return elem.prev
}

// NextLevel returns next element at specific level.
// If level is invalid, returns nil.
func (elem *Element[K, V]) NextLevel(level int) *Element[K, V] {
	if level < 0 || level >= len(elem.next) {
		return nil
	}

	return elem.next[level]
}

// Key returns the key of the elem.
func (elem *Element[K, V]) Key() K {
	return elem.key
}

// Level returns the level of this elem.
func (elem *Element[K, V]) Level() int {
	return len(elem.next)
}

func (elem *Element[K, V]) Index() int {
	return elem.list.Index(elem)
}
