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
				return &Element[K, V]{
					elementHeader: elementHeader[K, V]{
						levels: make([]*Element[K, V], 0, DefaultMaxLevel/2),
					},
				}
			},
		},
	}
}

//	func newElement[K, V any](list *SkipList[K, V], level int, key K, value V) *Element[K, V] {
//		return &Element[K, V]{
//			elementHeader: elementHeader[K, V]{
//				levels: make([]*Element[K, V], level),
//			},
//			Value: value,
//			key:   key,
//			list:  list,
//		}
//	}
func (f *elementPool[K, V]) Get(list *SkipList[K, V], level int, key K, value V) (element *Element[K, V]) {
	element = f.pool.Get().(*Element[K, V])
	for len(element.levels) < level {
		element.levels = append(element.levels, nil)
	}
	element.Value = value
	element.key = key
	element.list = list
	return
}

func (f *elementPool[K, V]) Put(element *Element[K, V]) {
	f.pool.Put(element)
	return
}
