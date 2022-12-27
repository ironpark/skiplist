// Copyright 2011 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package skiplist

//func TestCompareTypes(t *testing.T) {
//	a := assert.New(t)
//	cases := []struct {
//		kt       Comparable[int]
//		lhs, rhs interface{}
//		result   int
//	}{
//		{NumberComparator[int], 0, 0, 0},
//		{NumberComparator[int], 2, 0, 1},
//		{NumberComparator[int], -1, 1, -1},
//		{NumberComparator[byte], 9, 2, 1},
//		{NumberComparator[float64], 1.2, 1.20001, -1},
//		{BytesComparator[string], "foo", "bar", 1},
//		{BytesComparator[string], "001", "101", -1},
//		{BytesComparator[string], "equals", "equals", 0},
//		{BytesComparator[[]byte], []byte("abcdefghijk"), []byte("abcdefghij"), 1},
//	}
//
//	for i, c := range cases {
//		a.Use(&i, &c)
//		a.Equal(c.result, c.kt(c.lhs, c.rhs))
//	}
//}
