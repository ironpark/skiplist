package skiplist

import (
	"sync"
)

type elementPool[K, V any] struct {
	pool sync.Pool
}

func newElementPool[K, V any]() *elementPool[K, V] {
	return &elementPool[K, V]{
		pool: sync.Pool{
			New: func() interface{} {
				return &elementHeader[K, V]{
					next: make([]*Element[K, V], 0, DefaultMaxLevel/2),
				}
			},
		},
	}
}

func (f *elementPool[K, V]) Get(list *SkipList[K, V], level int, key K, value V) (element *Element[K, V]) {
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
