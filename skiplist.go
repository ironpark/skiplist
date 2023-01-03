// Copyright 2011 Huan Du. All rights reserved.
// Copyright 2022 Iron Park.

// Licensed under the MIT license that can be found in the LICENSE file.

// Package skiplist implement skip list data structure.
// See wikipedia for more details about this data structure. http://en.wikipedia.org/wiki/Skip_list
//
// Skip list is basically an ordered map.
//
// Here is a sample to use this package.
//
//	// Creates a new skip list and restricts key type to int-like types.
//	list := skiplist.New(skiplist.Int)
//
//	// Adds some values for keys.
//	list.Set(20, "Hello")
//	list.Set(10, "World")
//	list.Set(40, true) // Value type is not restricted.
//	list.Set(40, 1000) // Replace the of an existing element.
//
//	// Finds elements.
//	e := list.Get(10)           // Returns the element with the key.
//	_ = e.Value.(string)
//	v, ok := list.GetValue(20)  // Directly get value of the element. If the key is not found, ok is false.
//	v2 := list.MustGetValue(10) // Directly get value of the element. Panic if the key is not found.
//	notFound := list.Get(15)    // Returns nil if the key is not found.
//
//	// Removes an element and gets removed element.
//	old := list.Remove(40)
//	notFound := list.Remove(-20) // Returns nil if the key is not found.
//
//	// Initializes the list again to clean up all elements in the list.
//	list.Init()
package skiplist

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// DefaultMaxLevel is the default level for all newly created skip lists.
// It can be changed globally. Changing it will not affect existing lists.
// And all skip lists can update max level after creation through `SetMaxLevel()` method.
const (
	DefaultMaxLevel            = 18
	DefaultProbability float64 = 1 / math.E
)

// preallocDefaultMaxLevel is a constant to alloc memory on stack when Set new element.
const preallocDefaultMaxLevel = 48

type SkipList[K, V any] interface {
	Init() SkipList[K, V]
	SetProbability(newProbability float64)
	Front() (front *Element[K, V])
	Back() *Element[K, V]
	Len() int
	Set(key K, value V) (element *Element[K, V])
	FindNext(start *Element[K, V], key K) (next *Element[K, V])
	Find(key K) (elem *Element[K, V])
	Get(key K) (elem *Element[K, V])
	GetValue(key K) (val V, ok bool)
	MustGetValue(key K) V
	Remove(key K) (elem *Element[K, V])
	RemoveFront() (front *Element[K, V])
	RemoveBack() (back *Element[K, V])
	RemoveElement(elem *Element[K, V])
	MaxLevel() int
	Values() (values []V)
	Index(elem *Element[K, V]) (i int)
	Keys() (keys []K)
	SetMaxLevel(level int) (old int)
}

var _ = SkipList[int, int](&skipListUnSafe[int, int]{})

type skipListUnSafe[K, V any] struct {
	elementHeader[K, V]
	probTable      []float64
	pool           pool[K, V]
	comparable     Comparable[K]
	prevNodesCache []*elementHeader[K, V]
	rand           *rand.Rand

	maxLevel int
	length   int
	back     *Element[K, V]
}

// New creates a new skip list with comparable to compare keys.
//
// There are lots of pre-defined strict-typed keys like Int, Float64, String, etc.
// We can create custom comparable by implementing Comparable interface.
func New[K, V any](comparable Comparable[K], options ...Option) (skipList SkipList[K, V]) {
	option := &Options{
		maxLevel:    DefaultMaxLevel,
		probability: DefaultProbability,
		useLock:     false,
		usePool:     false,
	}
	for _, o := range options {
		o(option)
	}
	source := rand.NewSource(time.Now().UnixNano())
	sk := &skipListUnSafe[K, V]{
		elementHeader: elementHeader[K, V]{
			next: make([]*Element[K, V], option.maxLevel),
		},
		prevNodesCache: make([]*elementHeader[K, V], option.maxLevel),
		pool:           newElementPool[K, V](),
		probTable:      probabilityTable(option.probability, option.maxLevel),
		comparable:     comparable,
		rand:           rand.New(source),
		maxLevel:       DefaultMaxLevel,
	}
	if option.usePool {
		sk.pool = newElementPool[K, V]()
	} else {
		sk.pool = newFakePool[K, V]()
	}
	if option.useLock {
		return &safeSkipList[K, V]{
			skipListUnSafe: sk,
		}
	}
	return sk
}

// Init resets the list and discards all existing elements.
func (list *skipListUnSafe[K, V]) Init() SkipList[K, V] {
	list.back = nil
	list.length = 0
	list.next = make([]*Element[K, V], len(list.next))
	return list
}

// SetRandSource sets a new rand source.
//
// Skiplist uses global rand defined in math/rand by default.
// The default rand acquires a global mutex before generating any number.
// It's not necessary if the skiplist is well protected by caller.
func (list *skipListUnSafe[K, V]) SetRandSource(source rand.Source) {
	list.rand = rand.New(source)
}

// SetProbability changes the current P value of the list.
// It doesn't alter any existing data, only changes how future insert heights are calculated.
func (list *skipListUnSafe[K, V]) SetProbability(newProbability float64) {
	list.probTable = probabilityTable(newProbability, list.maxLevel)
}

// Front returns the first element.
//
// The complexity is O(1).
func (list *skipListUnSafe[K, V]) Front() (front *Element[K, V]) {
	return list.next[0]
}

// Back returns the last element.
//
// The complexity is O(1).
func (list *skipListUnSafe[K, V]) Back() *Element[K, V] {
	return list.back
}

// Len returns element count in this list.
//
// The complexity is O(1).
func (list *skipListUnSafe[K, V]) Len() int {
	return list.length
}

// Set sets value for the key.
// If the key exists, updates element's value.
// Returns the element holding the key and value.
//
// The complexity is O(log(N)).
func (list *skipListUnSafe[K, V]) Set(key K, value V) (element *Element[K, V]) {
	prevs := list.getPrevElementNodes(key)
	// replace
	if element = prevs[0].next[0]; element != nil && list.comparable(element.key, key) <= 0 {
		element.Value = value
		return element
	}
	// insert
	nextElement := prevs[0].next[0]
	element = list.pool.Get(list, list.randLevel(), key, value)

	for i := range element.next {
		element.next[i] = prevs[i].next[i]
		prevs[i].next[i] = element
	}
	if nextElement == nil {
		// 끝에 인서트
		element.prev = list.back
		list.back = element
	} else {
		// 앞 혹은 중간
		element.prev = nextElement.prev
		nextElement.prev = element
	}
	list.length++
	return
}

// FindNext returns the first element after start that is greater or equal to key.
// If start is greater or equal to key, returns start.
// If there is no such element, returns nil.
// If start is nil, find element from front.
//
// The complexity is O(log(N)).
func (list *skipListUnSafe[K, V]) FindNext(start *Element[K, V], key K) (next *Element[K, V]) {
	if list.length == 0 {
		return
	}
	var header = &list.elementHeader
	maxLevel := list.maxLevel
	if start != nil {
		if list.comparable(key, start.key) <= 0 {
			return start
		}
		header = start.elementHeader
		maxLevel = start.Level()

	}
	if list.comparable(key, list.Front().key) <= 0 {
		return list.Front()
	}
	if list.comparable(key, list.Back().key) > 0 {
		return
	}

	for i := maxLevel - 1; i >= 0; i-- {
		next = header.next[i]
		// 입력키가 다음키보다 크면 점프
		for next != nil {
			switch list.comparable(key, next.key) {
			case 0:
				// key == next.key
				return next
			case 1:
				// key > next.key
				header = next.elementHeader
				next = next.next[i]
			case -1:
				// key < next.key
				goto Next
			}
		}
	Next:
	}
	return next
}

// Find returns the first element that is greater or equal to key.
// It's short hand for FindNext(nil, key).
//
// The complexity is O(log(N)).
func (list *skipListUnSafe[K, V]) Find(key K) (elem *Element[K, V]) {
	return list.FindNext(nil, key)
}

// Get returns an element with the key.
// If the key is not found, returns nil.
//
// The complexity is O(log(N)).
func (list *skipListUnSafe[K, V]) Get(key K) (elem *Element[K, V]) {
	var prev = &list.elementHeader
	var next *Element[K, V]

	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && list.comparable(key, next.key) > 0 {
			prev = next.elementHeader
			next = next.next[i]
		}
	}

	if next != nil && list.comparable(next.key, key) <= 0 {
		return next
	}
	return
}

// GetValue returns value of the element with the key.
// It's short hand for Get().Value.
//
// The complexity is O(log(N)).
func (list *skipListUnSafe[K, V]) GetValue(key K) (val V, ok bool) {
	element := list.Get(key)
	if element == nil {
		return
	}
	val = element.Value
	ok = true
	return
}

// MustGetValue returns value of the element with the key.
// It will panic if the key doesn't exist in the list.
//
// The complexity is O(log(N)).
func (list *skipListUnSafe[K, V]) MustGetValue(key K) V {
	element := list.Get(key)
	if element == nil {
		panic(fmt.Errorf("skiplist: cannot find key `%v` in skiplist", key))
	}
	return element.Value
}

// Remove removes an element.
// Returns removed element pointer if found, nil if it's not found.
//
// The complexity is O(log(N)).
func (list *skipListUnSafe[K, V]) Remove(key K) (elem *Element[K, V]) {
	prevs := list.getPrevElementNodes(key)
	elem = prevs[0].next[0]
	if elem == nil {
		return nil
	}
	if list.comparable(elem.key, key) != 0 {
		return nil
	}
	tail := elem.next[0] == nil
	if elem.next[0] != nil && elem.next[0].prev != nil {
		elem.next[0].prev = elem.prev
	}
	for k, v := range elem.next {
		prevs[k].next[k] = v
	}
	if tail {
		list.back = elem.prev
	}
	list.length--
	list.pool.Put(elem)
	return
}

// RemoveFront removes front element node and returns the removed element.
//
// The complexity is O(1).
func (list *skipListUnSafe[K, V]) RemoveFront() (front *Element[K, V]) {
	if list.length == 0 {
		return
	}
	front = list.Front()
	list.RemoveElement(front)
	return
}

// RemoveBack removes back element node and returns the removed element.
//
// The complexity is O(log(N)).
func (list *skipListUnSafe[K, V]) RemoveBack() (back *Element[K, V]) {
	if list.length == 0 {
		return
	}
	back = list.back
	list.RemoveElement(back)
	return
}

// RemoveElement removes the elem from the list.
//
// The complexity is O(log(N)).
func (list *skipListUnSafe[K, V]) RemoveElement(elem *Element[K, V]) {
	if elem == nil || elem.list != list {
		return
	}
	_ = list.Remove(elem.key)
}

// MaxLevel returns current max level value.
func (list *skipListUnSafe[K, V]) MaxLevel() int {
	return list.maxLevel
}

// Values returns list of values
func (list *skipListUnSafe[K, V]) Values() (values []V) {
	for el := list.Front(); el != nil; el = el.Next() {
		values = append(values, el.Value)
	}
	return
}

// Index returns index of element
func (list *skipListUnSafe[K, V]) Index(elem *Element[K, V]) (i int) {
	if elem == nil {
		return -1
	}
	for e := elem.Prev(); e != nil; e = e.Prev() {
		i++
	}
	return i
}

// Keys returns list of keys
func (list *skipListUnSafe[K, V]) Keys() (keys []K) {
	for el := list.Front(); el != nil; el = el.Next() {
		keys = append(keys, el.key)
	}
	return
}

// SetMaxLevel changes skip list max level.
// If level is not greater than 0, just panic.
func (list *skipListUnSafe[K, V]) SetMaxLevel(level int) (old int) {
	if level <= 0 {
		panic(fmt.Errorf("skiplist: level must be larger than 0 (current is %v)", level))

	}
	for len(list.prevNodesCache) < level {
		list.prevNodesCache = append(list.prevNodesCache, nil)
	}
	list.probTable = probabilityTable(DefaultProbability, level)
	list.maxLevel = level
	old = len(list.next)

	if level == old {
		return
	}

	if old > level {
		for i := old - 1; i >= level; i-- {
			if list.next[i] != nil {
				level = i
				break
			}
		}
		list.next = list.next[:level]
		return
	}

	if level <= cap(list.next) {
		list.next = list.next[:level]
		return
	}

	levels := make([]*Element[K, V], level)
	copy(levels, list.next)
	list.next = levels
	return
}

func (list *skipListUnSafe[K, V]) randLevel() (level int) {
	r := float64(list.rand.Int63()) / (1 << 63)
	for level = 1; level < list.maxLevel && r < list.probTable[level]; level++ {

	}
	return
}

// getPrevElementNodes is the private search mechanism that other functions use.
// Finds the previous nodes on each level relative to the current Element and
// caches them. This approach is similar to a "search finger" as described by Pugh:
// http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.17.524
// original by https://github.com/sean-public/fast-skiplist
func (list *skipListUnSafe[K, V]) getPrevElementNodes(key K) (prevs []*elementHeader[K, V]) {
	prev := &list.elementHeader
	prevs = list.prevNodesCache
	for i := list.maxLevel - 1; i >= 0; i-- {
		next := prev.next[i]
		for next != nil && list.comparable(key, next.key) > 0 {
			prev = next.elementHeader
			next = next.next[i]
		}
		prevs[i] = prev
	}
	return
}

// probabilityTable calculates in advance the probability of a new node having a given level.
// probability is in [0, 1], MaxLevel is [0, 64]
// Returns a table of floating point probabilities that each level should be included during an insert.
// original by https://github.com/sean-public/fast-skiplist
func probabilityTable(probability float64, MaxLevel int) (table []float64) {
	for i := 1; i <= MaxLevel; i++ {
		prob := math.Pow(probability, float64(i-1))
		table = append(table, prob)
	}

	return table
}
