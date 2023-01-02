# Skip List in Golang

[![Go](https://github.com/ironpark/skiplist/workflows/Go/badge.svg)](https://github.com/ironpark/skiplist/actions)
[![Go Doc](https://godoc.org/github.com/ironpark/skiplist?status.svg)](https://pkg.go.dev/github.com/ironpark/skiplist)
[![Go Report](https://goreportcard.com/badge/github.com/Ironpark/skiplist)](https://goreportcard.com/report/github.com/Ironpark/skiplist)
[![Coverage Status](https://coveralls.io/repos/github/Ironpark/skiplist/badge.svg?branch=master)](https://coveralls.io/github/Ironpark/skiplist?branch=master)

This package was created by forking the skiplist library that [Huandu's](https://github.com/huandu) great job

Skip list is an ordered map. See wikipedia page [skip list](http://en.wikipedia.org/wiki/Skip_list) to learn algorithm details about this data structure.

Highlights in this implementation:
- Based on generic  **No reflection, No interface{}**
- Support custom comparable function so that any type can be used as key.
- Rand source and max level can be changed per list. It can be useful in performance critical scenarios.
- Optional memory pool use
- Optional thread-safe instance `SafeSkipList`, `SkipList`

## Warrning
Not fully tested

## Install

Install this package through `go get`.

```bash
go get github.com/ironpark/skiplist
```
## Example

### Basic Usage

Here is a quick sample.

```go
package main

import (
    "fmt"
    "github.com/ironpark/skiplist"
)

func main() {
    // Create a skip list with int key.
    list := skiplist.New[int,any](skiplist.NumberComparator[int])

    // Add some values. Value can be anything.
    list.Set(12, "hello world")
    list.Set(34, 56)
    list.Set(78, 90.12)

    // Get element by index.
    elem := list.Get(34)                // Value is stored in elem.Value.
    fmt.Println(elem.Value)             // Output: 56
    next := elem.Next()                 // Get next element.
    prev := next.Prev()                 // Get previous element.
    fmt.Println(next.Value, prev.Value) // Output: 90.12    56

    // Or, directly get value just like a map
    val, ok := list.GetValue(34)
    fmt.Println(val, ok) // Output: 56  true

    // Find first elements with score greater or equal to key
    foundElem := list.Find(30)
    fmt.Println(foundElem.Key(), foundElem.Value) // Output: 34 56

    // Remove an element for key.
    list.Remove(34)
}
```

### Thread Safe

```go
package main

import (
  "fmt"
  "github.com/ironpark/skiplist"
  "sync"
)

func main() {
  // Create a skip list with int key.
  list := skiplist.New[int, struct{}](skiplist.NumberComparator[int],skiplist.WithMutex())
  wg := sync.WaitGroup{}
  wg.Add(100)
  for i := 0; i < 100; i++ {
    go func(i int) {
      list.Set(i, struct{}{})
      wg.Done()
    }(i)
  }
  wg.Wait()
  fmt.Println(list.Keys())
}
```
## License

This library is licensed under MIT license. See LICENSE for details.
