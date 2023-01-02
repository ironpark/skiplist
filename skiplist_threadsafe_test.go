package skiplist

import (
	"sync"
	"testing"
)

func TestSafeSkipList_Set(t *testing.T) {
	list := New[int, struct{}](NumberComparator[int], WithMutex())
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			list.Set(i, struct{}{})
			wg.Done()
		}(i)
	}
	wg.Wait()
	//fmt.Println(list.Keys())
}
