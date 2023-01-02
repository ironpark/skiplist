package skiplist

import (
	"sync"
)

type pool[K, V any] interface {
	Get(list SkipList[K, V], level int, key K, value V) (element *Element[K, V])
	Put(element *Element[K, V])
}
type elementPool[K, V any] struct {
	pool sync.Pool
}

func newElementPool[K, V any]() *elementPool[K, V] {
	return &elementPool[K, V]{
		pool: sync.Pool{
			New: func() interface{} {
				return &elementHeader[K, V]{
					next: make([]*Element[K, V], 0, DefaultMaxLevel),
				}
			},
		},
	}
}

func (f *elementPool[K, V]) Get(list SkipList[K, V], level int, key K, value V) (element *Element[K, V]) {
	header := f.pool.Get().(*elementHeader[K, V])
	header.next = header.next[:level]
	return &Element[K, V]{
		list:          list,
		Value:         value,
		key:           key,
		elementHeader: header,
	}
}

func (f *elementPool[K, V]) Put(element *Element[K, V]) {
	element.list = nil
	element.prev = nil
	next := element.next
	for i := range next {
		next[i] = nil
	}
	f.pool.Put(element.elementHeader)
	return
}

type fakePool[K, V any] struct {
}

func newFakePool[K, V any]() *fakePool[K, V] {
	return &fakePool[K, V]{}
}

func (f *fakePool[K, V]) Get(list SkipList[K, V], level int, key K, value V) (element *Element[K, V]) {
	return &Element[K, V]{
		list:  list,
		Value: value,
		key:   key,
		elementHeader: &elementHeader[K, V]{
			make([]*Element[K, V], level),
		},
	}
}

func (f *fakePool[K, V]) Put(element *Element[K, V]) {
	element.list = nil
	element.prev = nil
	element.next = nil
	return
}
