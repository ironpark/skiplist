package skiplist

import "golang.org/x/exp/constraints"

type Numbers interface {
	constraints.Integer | constraints.Float | rune
}

type Bytes interface {
	~[]byte | ~string
}

func NumberComparator[K Numbers](lk, rk K) int {
	if lk > rk {
		return 1
	}
	if lk < rk {
		return -1
	}
	return 0
}

func BytesComparator[K Bytes](lk, rk K) int {
	lhs, rhs := bytesScore(lk), bytesScore(rk)
	if lhs > rhs {
		return 1
	}
	if lhs < rhs {
		return -1
	}
	return 0
}

func bytesScore[K Bytes](data K) (score uint64) {
	l := len(data)
	// only use first 8 bytes
	if l > 8 {
		l = 8
	}
	// Consider str as a Big-Endian uint64.
	for i := 0; i < l; i++ {
		shift := uint(64 - 8 - i*8)
		score |= uint64(data[i]) << shift
	}
	return
}
