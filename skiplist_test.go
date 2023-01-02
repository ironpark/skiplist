// Copyright 2011 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package skiplist

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestSkipList_Set(t *testing.T) {
	a := assert.New(t)
	list := New[int, int](NumberComparator[int])
	var keys, values []int
	var keysShuffled, valuesShuffled []int
	for i := 0; i < 100; i++ {
		keys = append(keys, i)
		values = append(values, i*2)
		keysShuffled = append(keysShuffled, i)
		valuesShuffled = append(valuesShuffled, i*2)
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keysShuffled[i], keysShuffled[j] = keysShuffled[j], keysShuffled[i]
		valuesShuffled[i], valuesShuffled[j] = valuesShuffled[j], valuesShuffled[i]
	})

	for i := 0; i < len(keys); i++ {
		list.Set(keysShuffled[i], valuesShuffled[i])
	}
	a.Equal(list.Values(), values)
	a.Equal(list.Keys(), keys)
}

func TestInsertAndRemove(t *testing.T) {
	list := New[int, any](NumberComparator[int])
	list.Set(1, "#1")
	list.Remove(1)
	list.Set(1, "#1")
	list.Remove(3)
	list.Set(1, "#1")
	list.Set(2, "#2")
	list.Set(3, "#3")
	list.Set(4, "#4")
	list.Remove(3)
}

func TestDuplicate(t *testing.T) {
	a := assert.New(t)
	list := New[int, any](NumberComparator[int])
	list.Set(1, "#1")
	list.Set(1, "#4")
	list.Set(1, "#2")
	list.Set(1, "#3")
	a.True(list.Len() == 1)
	a.Equal(list.Front(), list.Back())
	a.Equal(list.Front().Value, "#3")
}

func TestIndex(t *testing.T) {
	a := assert.New(t)
	list := New[int, any](NumberComparator[int])
	list.Set(0, "0")
	list.Set(1, "1")
	list.Set(2, "2")
	a.Equal(list.Index(list.Get(0)), 0)
	a.Equal(list.Index(list.Get(1)), 1)
	a.Equal(list.Index(list.Get(2)), 2)
	list.Remove(1)

	a.Equal(list.Get(0).Index(), 0)
	a.Equal(list.Get(2).Index(), 1)
}

func TestSkipList_FindNext(t *testing.T) {
	a := assert.New(t)
	list := New[float64, any](NumberComparator[float64])
	list.Set(0, "0")
	a.True(list.FindNext(nil, 0) != nil)
	list.Set(1, "1")
	a.True(list.FindNext(nil, 0) != nil)
	a.True(list.FindNext(nil, 1) != nil)
	list.Set(2, "2")
	a.True(list.FindNext(nil, 0) != nil)
	a.True(list.FindNext(nil, 1) != nil)
	a.True(list.FindNext(nil, 2) != nil)
	a.True(list.FindNext(nil, 3) == nil)
}

func TestBasicCRUD(t *testing.T) {
	a := assert.New(t)
	list := New[float64, any](NumberComparator[float64])
	a.True(list.Len() == 0)
	a.True(list.Find(0) == nil)

	elem1 := list.Set(12.34, "first")
	a.True(elem1 != nil)
	a.Equal(list.Len(), 1)
	a.Equal(list.Front(), elem1)
	a.Equal(list.Back(), elem1)
	a.True(elem1.Next() == nil)
	a.True(elem1.Prev() == nil)
	a.Equal(list.Find(0).key, elem1.key)
	a.Equal(list.Find(12.34), elem1)
	a.True(list.Find(15) == nil)

	//assertSanity(a, list)
	//
	elem2 := list.Set(23.45, "second")
	a.True(elem2 != nil)
	a.NotEqual(elem1, elem2)
	a.Equal(list.Len(), 2)
	a.Equal(list.Front(), elem1)
	a.Equal(list.Back(), elem2)
	a.True(elem2.Next() == nil)
	a.Equal(elem2.Prev(), elem1)
	a.Equal(list.Find(-10), elem1)
	a.Equal(list.Find(15), elem2)
	a.True(list.Find(25) == nil)

	//assertSanity(a, list)
	//
	elem3 := list.Set(16.78, "middle")
	a.True(elem3 != nil)
	a.NotEqual(elem3, elem1)
	a.NotEqual(elem3, elem2)
	a.Equal(list.Len(), 3)
	a.Equal(list.Front(), elem1)
	a.Equal(list.Back(), elem2)
	a.Equal(elem3.Next(), elem2)
	a.Equal(elem3.Prev(), elem1)
	a.Equal(list.Find(-20), elem1)
	a.Equal(list.Find(15), elem3)
	a.Equal(list.Find(20), elem2)

	//assertSanity(a, list)
	//
	elem4 := list.Set(9.01, "very beginning")
	a.True(elem4 != nil)
	a.NotEqual(elem4, elem1)
	a.NotEqual(elem4, elem2)
	a.NotEqual(elem4, elem3)
	a.Equal(list.Len(), 4)
	a.Equal(list.Front(), elem4)
	a.Equal(list.Back(), elem2)
	a.Equal(elem4.Next(), elem1)
	a.True(elem4.Prev() == nil)
	a.Equal(list.Find(0), elem4)
	a.Equal(list.Find(15), elem3)
	a.Equal(list.Find(20), elem2)
	//
	////assertSanity(a, list)
	//
	elem5 := list.Set(16.78, "middle overwrite")
	a.True(elem3 != nil)
	a.NotEqual(elem3, elem1)
	a.NotEqual(elem3, elem2)
	a.Equal(elem5, elem3)
	a.NotEqual(elem5, elem4)
	a.Equal(list.Len(), 4)
	a.Equal(list.Front(), elem4)
	a.Equal(list.Back(), elem2)
	a.Equal(elem5.Next(), elem2)
	a.Equal(elem5.Prev(), elem1)
	a.Equal(list.Find(15), elem5)
	a.Equal(list.Find(16.78), elem5)
	a.Equal(list.Find(16.79), elem2)
	a.Equal(list.FindNext(nil, 15), elem5)
	a.Equal(list.FindNext(nil, 16.78), elem5)
	a.Equal(list.FindNext(nil, 16.79), elem2)
	a.Equal(list.FindNext(elem1, 15), elem5)
	a.Equal(list.FindNext(elem5, 15), elem5)
	a.True(list.FindNext(elem5, 30) == nil)
	//
	min1_2 := func(a, b int) int {
		if a < b {
			return a / 2
		}
		return b / 2
	}
	a.Equal(elem5.NextLevel(0), elem5.Next())
	a.True(elem5.NextLevel(-1) == nil)
	a.Equal(elem5.NextLevel(min1_2(elem2.Level(), elem5.Level())), elem2)
	a.True(elem5.NextLevel(elem2.Level()) == nil)
	//a.Equal(elem5.PrevLevel(0), elem5.Prev())
	//a.Equal(elem5.PrevLevel(min1_2(elem1.Level(), elem5.Level())), elem1)
	//a.Equal(elem5.PrevLevel(-1), nil)

	a.True(list.Remove(9999) == nil)
	a.Equal(list.Len(), 4)
	a.True(list.Remove(13.24) == nil)
	a.Equal(list.Len(), 4)

	list.SetMaxLevel(1)
	list.SetMaxLevel(128)
	list.SetMaxLevel(32)
	list.SetMaxLevel(32)

	elem2Removed := list.Remove(elem2.Key())
	a.True(elem2Removed != nil)
	a.Equal(elem2Removed, elem2)
	a.True(elem2Removed.Prev() == nil)
	a.True(elem2Removed.Next() == nil)
	a.Equal(list.Len(), 3)
	a.Equal(list.Front(), elem4)
	a.Equal(list.Back(), elem5)
	a.Equal(list.Find(-99), elem4)
	a.Equal(list.Find(10), elem1)
	a.Equal(list.Find(15), elem3)
	a.True(list.Find(20) == nil)

	front := list.RemoveFront()
	a.True(front == elem4)
	a.Equal(list.Len(), 2)
	a.Equal(list.Front(), elem1)
	a.Equal(list.Back(), elem5)
	a.Equal(list.Find(-99), elem1)

	back := list.RemoveBack()
	a.True(back == elem5)
	a.Equal(list.Len(), 1)
	a.Equal(list.Front(), elem1)
	a.Equal(list.Back(), elem1)
	a.True(list.Find(15) == nil)
	a.Equal(list.FindNext(nil, 10), elem1)
	a.Equal(list.FindNext(elem1, 10), elem1)
	a.True(list.FindNext(nil, 15) == nil)

	list.Init()
	a.Equal(list.Len(), 0)
	a.True(list.Get(12.34) == nil)
}

//238001384
//333830583
