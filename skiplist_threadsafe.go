// Copyright 2022 Iron Park. All rights reserved.
package skiplist

import (
	"math/rand"
	"sync"
)

// SafeSkipList is the header of a skip list.
type SafeSkipList[K, V any] struct {
	*SkipList[K, V]
	lock sync.RWMutex
}

func NewSafe[K, V any](comparable Comparable[K]) *SafeSkipList[K, V] {
	return &SafeSkipList[K, V]{
		SkipList: New[K, V](comparable),
	}
}

// Init resets the list and discards all existing elements.
func (list *SafeSkipList[K, V]) Init() *SkipList[K, V] {
	list.lock.Lock()
	defer list.lock.Unlock()
	return list.SkipList.Init()
}

// SetRandSource sets a new rand source.
//
// Skiplist uses global rand defined in math/rand by default.
// The default rand acquires a global mutex before generating any number.
// It's not necessary if the skiplist is well protected by caller.
func (list *SafeSkipList[K, V]) SetRandSource(source rand.Source) {
	list.lock.Lock()
	defer list.lock.Unlock()
	list.SkipList.SetRandSource(source)
}

// SetProbability changes the current P value of the list.
// It doesn't alter any existing data, only changes how future insert heights are calculated.
func (list *SafeSkipList[K, V]) SetProbability(newProbability float64) {
	list.lock.Lock()
	defer list.lock.Unlock()
	list.SkipList.SetProbability(newProbability)
}

// Front returns the first element.
//
// The complexity is O(1).
func (list *SafeSkipList[K, V]) Front() (front *Element[K, V]) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.SkipList.Front()
}

// Back returns the last element.
//
// The complexity is O(1).
func (list *SafeSkipList[K, V]) Back() *Element[K, V] {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.SkipList.Back()
}

// Len returns element count in this list.
//
// The complexity is O(1).
func (list *SafeSkipList[K, V]) Len() int {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.SkipList.Len()
}

// Set sets value for the key.
// If the key exists, updates element's value.
// Returns the element holding the key and value.
//
// The complexity is O(log(N)).
func (list *SafeSkipList[K, V]) Set(key K, value V) (elem *Element[K, V]) {
	list.lock.Lock()
	defer list.lock.Unlock()
	return list.SkipList.Set(key, value)
}

func (list *SafeSkipList[K, V]) FindNext(start *Element[K, V], key K) (elem *Element[K, V]) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.SkipList.FindNext(start, key)
}

// Find returns the first element that is greater or equal to key.
// It's short hand for FindNext(nil, key).
//
// The complexity is O(log(N)).
func (list *SafeSkipList[K, V]) Find(key K) (elem *Element[K, V]) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.SkipList.Find(key)
}

// Get returns an element with the key.
// If the key is not found, returns nil.
//
// The complexity is O(log(N)).
func (list *SafeSkipList[K, V]) Get(key K) (elem *Element[K, V]) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.SkipList.Get(key)
}

// GetValue returns value of the element with the key.
// It's short hand for Get().Value.
//
// The complexity is O(log(N)).
func (list *SafeSkipList[K, V]) GetValue(key K) (val V, ok bool) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.SkipList.GetValue(key)
}

// MustGetValue returns value of the element with the key.
// It will panic if the key doesn't exist in the list.
//
// The complexity is O(log(N)).
func (list *SafeSkipList[K, V]) MustGetValue(key K) V {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.SkipList.MustGetValue(key)
}

// Remove removes an element.
// Returns removed element pointer if found, nil if it's not found.
//
// The complexity is O(log(N)).
func (list *SafeSkipList[K, V]) Remove(key K) (elem *Element[K, V]) {
	list.lock.Lock()
	defer list.lock.Unlock()
	return list.SkipList.Remove(key)
}

// RemoveFront removes front element node and returns the removed element.
//
// The complexity is O(1).
func (list *SafeSkipList[K, V]) RemoveFront() (front *Element[K, V]) {
	list.lock.Lock()
	defer list.lock.Unlock()
	return list.SkipList.RemoveFront()
}

// RemoveBack removes back element node and returns the removed element.
//
// The complexity is O(log(N)).
func (list *SafeSkipList[K, V]) RemoveBack() (back *Element[K, V]) {
	list.lock.Lock()
	defer list.lock.Unlock()
	return list.SkipList.RemoveBack()
}

// RemoveElement removes the elem from the list.
//
// The complexity is O(log(N)).
func (list *SafeSkipList[K, V]) RemoveElement(elem *Element[K, V]) {
	list.lock.Lock()
	defer list.lock.Unlock()
	list.SkipList.RemoveElement(elem)
}

// MaxLevel returns current max level value.
func (list *SafeSkipList[K, V]) MaxLevel() int {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.maxLevel
}

// Values returns list of values
func (list *SafeSkipList[K, V]) Values() (values []V) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.SkipList.Values()
}

// Keys returns list of keys
func (list *SafeSkipList[K, V]) Keys() (keys []K) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.SkipList.Keys()
}

// SetMaxLevel changes skip list max level.
// If level is not greater than 0, just panic.
func (list *SafeSkipList[K, V]) SetMaxLevel(level int) (old int) {
	list.lock.Lock()
	defer list.lock.Unlock()
	list.SkipList.SetMaxLevel(level)
	return
}
