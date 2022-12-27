// Copyright 2011 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package skiplist

import (
	"unsafe"
)

// Element is an element node of a skip list.
type Element[K, V any] struct {
	elementHeader[K, V]

	Value        V
	key          K
	prev         *Element[K, V]  // Points to previous adjacent elem.
	prevTopLevel *Element[K, V]  // Points to previous element which points to this element's top most level.
	list         *SkipList[K, V] // The list contains this elem.
}

// elementHeader is the header of an element or a skip list.
// It must be the first anonymous field in a type to make Element() work correctly.
type elementHeader[K, V any] struct {
	levels []*Element[K, V] // Next element at all levels.
}

func (header *elementHeader[K, V]) Element() *Element[K, V] {
	return (*Element[K, V])(unsafe.Pointer(header))
}

func newElement[K, V any](list *SkipList[K, V], level int, key K, value V) *Element[K, V] {
	return &Element[K, V]{
		elementHeader: elementHeader[K, V]{
			levels: make([]*Element[K, V], level),
		},
		Value: value,
		key:   key,
		list:  list,
	}
}

// Next returns next adjacent elem.
func (elem *Element[K, V]) Next() *Element[K, V] {
	if len(elem.levels) == 0 {
		return nil
	}
	return elem.levels[0]
}

// Prev returns previous adjacent elem.
func (elem *Element[K, V]) Prev() *Element[K, V] {
	return elem.prev
}

// NextLevel returns next element at specific level.
// If level is invalid, returns nil.
func (elem *Element[K, V]) NextLevel(level int) *Element[K, V] {
	if level < 0 || level >= len(elem.levels) {
		return nil
	}

	return elem.levels[level]
}

// PrevLevel returns previous element which points to this element at specific level.
// If level is invalid, returns nil.
func (elem *Element[K, V]) PrevLevel(level int) *Element[K, V] {
	if level < 0 || level >= len(elem.levels) {
		return nil
	}

	if level == 0 {
		return elem.prev
	}

	if level == len(elem.levels)-1 {
		return elem.prevTopLevel
	}

	prev := elem.prev

	for prev != nil {
		if level < len(prev.levels) {
			return prev
		}

		prev = prev.prevTopLevel
	}

	return prev
}

// Key returns the key of the elem.
func (elem *Element[K, V]) Key() K {
	return elem.key
}

// Level returns the level of this elem.
func (elem *Element[K, V]) Level() int {
	return len(elem.levels)
}

func (elem *Element[K, V]) reset() {
	elem.list = nil
	elem.prev = nil
	elem.prevTopLevel = nil
	elem.levels = elem.levels[:0]
}

func resetLevels[K, V any](levels []*Element[K, V]) {
	levels[0] = nil
	for bp := 1; bp < len(levels); bp *= 2 {
		copy(levels[bp:], levels[:bp])
	}
}
