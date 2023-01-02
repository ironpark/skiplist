// Copyright 2022 Iron Park. All rights reserved.
package skiplist

import (
	"math/rand"
	"sync"
)

// SafeSkipList is the header of a skip list.
type safeSkipList[K, V any] struct {
	*skipListUnSafe[K, V]
	lock sync.RWMutex
}

// Init resets the list and discards all existing elements.
func (list *safeSkipList[K, V]) Init() SkipList[K, V] {
	list.lock.Lock()
	defer list.lock.Unlock()
	return list.skipListUnSafe.Init()
}

// SetRandSource sets a new rand source.
//
// Skiplist uses global rand defined in math/rand by default.
// The default rand acquires a global mutex before generating any number.
// It's not necessary if the skiplist is well protected by caller.
func (list *safeSkipList[K, V]) SetRandSource(source rand.Source) {
	list.lock.Lock()
	defer list.lock.Unlock()
	list.skipListUnSafe.SetRandSource(source)
}

// SetProbability changes the current P value of the list.
// It doesn't alter any existing data, only changes how future insert heights are calculated.
func (list *safeSkipList[K, V]) SetProbability(newProbability float64) {
	list.lock.Lock()
	defer list.lock.Unlock()
	list.skipListUnSafe.SetProbability(newProbability)
}

// Front returns the first element.
//
// The complexity is O(1).
func (list *safeSkipList[K, V]) Front() (front *Element[K, V]) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.skipListUnSafe.Front()
}

// Back returns the last element.
//
// The complexity is O(1).
func (list *safeSkipList[K, V]) Back() *Element[K, V] {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.skipListUnSafe.Back()
}

// Len returns element count in this list.
//
// The complexity is O(1).
func (list *safeSkipList[K, V]) Len() int {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.skipListUnSafe.Len()
}

// Set sets value for the key.
// If the key exists, updates element's value.
// Returns the element holding the key and value.
//
// The complexity is O(log(N)).
func (list *safeSkipList[K, V]) Set(key K, value V) (elem *Element[K, V]) {
	list.lock.Lock()
	defer list.lock.Unlock()
	return list.skipListUnSafe.Set(key, value)
}

func (list *safeSkipList[K, V]) FindNext(start *Element[K, V], key K) (elem *Element[K, V]) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.skipListUnSafe.FindNext(start, key)
}

// Find returns the first element that is greater or equal to key.
// It's short hand for FindNext(nil, key).
//
// The complexity is O(log(N)).
func (list *safeSkipList[K, V]) Find(key K) (elem *Element[K, V]) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.skipListUnSafe.Find(key)
}

// Get returns an element with the key.
// If the key is not found, returns nil.
//
// The complexity is O(log(N)).
func (list *safeSkipList[K, V]) Get(key K) (elem *Element[K, V]) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.skipListUnSafe.Get(key)
}

// GetValue returns value of the element with the key.
// It's short hand for Get().Value.
//
// The complexity is O(log(N)).
func (list *safeSkipList[K, V]) GetValue(key K) (val V, ok bool) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.skipListUnSafe.GetValue(key)
}

// MustGetValue returns value of the element with the key.
// It will panic if the key doesn't exist in the list.
//
// The complexity is O(log(N)).
func (list *safeSkipList[K, V]) MustGetValue(key K) V {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.skipListUnSafe.MustGetValue(key)
}

// Remove removes an element.
// Returns removed element pointer if found, nil if it's not found.
//
// The complexity is O(log(N)).
func (list *safeSkipList[K, V]) Remove(key K) (elem *Element[K, V]) {
	list.lock.Lock()
	defer list.lock.Unlock()
	return list.skipListUnSafe.Remove(key)
}

// RemoveFront removes front element node and returns the removed element.
//
// The complexity is O(1).
func (list *safeSkipList[K, V]) RemoveFront() (front *Element[K, V]) {
	list.lock.Lock()
	defer list.lock.Unlock()
	return list.skipListUnSafe.RemoveFront()
}

// RemoveBack removes back element node and returns the removed element.
//
// The complexity is O(log(N)).
func (list *safeSkipList[K, V]) RemoveBack() (back *Element[K, V]) {
	list.lock.Lock()
	defer list.lock.Unlock()
	return list.skipListUnSafe.RemoveBack()
}

// RemoveElement removes the elem from the list.
//
// The complexity is O(log(N)).
func (list *safeSkipList[K, V]) RemoveElement(elem *Element[K, V]) {
	list.lock.Lock()
	defer list.lock.Unlock()
	list.skipListUnSafe.RemoveElement(elem)
}

// MaxLevel returns current max level value.
func (list *safeSkipList[K, V]) MaxLevel() int {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.maxLevel
}

// Index returns index of element
func (list *safeSkipList[K, V]) Index(elem *Element[K, V]) (i int) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.skipListUnSafe.Index(elem)
}

// Values returns list of values
func (list *safeSkipList[K, V]) Values() (values []V) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.skipListUnSafe.Values()
}

// Keys returns list of keys
func (list *safeSkipList[K, V]) Keys() (keys []K) {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.skipListUnSafe.Keys()
}

// SetMaxLevel changes skip list max level.
// If level is not greater than 0, just panic.
func (list *safeSkipList[K, V]) SetMaxLevel(level int) (old int) {
	list.lock.Lock()
	defer list.lock.Unlock()
	list.skipListUnSafe.SetMaxLevel(level)
	return
}
