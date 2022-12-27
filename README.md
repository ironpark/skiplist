# Skip List in Golang

[![Go](https://github.com/ironpark/skiplist/workflows/Go/badge.svg)](https://github.com/ironpark/skiplist/actions)
[![Go Doc](https://godoc.org/github.com/ironpark/skiplist?status.svg)](https://pkg.go.dev/github.com/huandu/skiplist)
[![Go Report](https://goreportcard.com/badge/github.com/ironpark/skiplist)](https://goreportcard.com/report/github.com/ironpark/skiplist)
[![Coverage Status](https://coveralls.io/repos/github/ironpark/skiplist/badge.svg?branch=master)](https://coveralls.io/github/ironpark/skiplist?branch=master)

This package was created by forking the skiplist library that [Huandu's](https://github.com/huandu) great job

Skip list is an ordered map. See wikipedia page [skip list](http://en.wikipedia.org/wiki/Skip_list) to learn algorithm details about this data structure.

Highlights in this implementation:

- Support custom comparable function so that any type can be used as key.
- Rand source and max level can be changed per list. It can be useful in performance critical scenarios.

## Install

Install this package through `go get`.

```bash
    go get github.com/ironpark/skiplist
```

## Basic Usage

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


## License

This library is licensed under MIT license. See LICENSE for details.
