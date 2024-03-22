package threadsave

import "sync/atomic"

type node[T any] struct {
	val  T
	next atomic.Pointer[node[T]]
}

func newNode[T any](val T, ptr *node[T]) *node[T] {
	nd := node[T]{val: val}
	nd.next.Store(ptr)
	return &nd
}
