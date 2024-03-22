package threadsave

import (
	"fmt"
	"parallel-prog/1/stack"
	"sync/atomic"
)

type node[T any] struct {
	val  T
	next atomic.Pointer[node[T]]
}

type CommonStack[T any] struct {
	top atomic.Pointer[node[T]]
}

func newNode[T any](val T, ptr *node[T]) *node[T] {
	nd := node[T]{val: val}
	nd.next.Store(ptr)
	return &nd
}

func NewCommonStack[T any]() *CommonStack[T] {
	return &CommonStack[T]{}
}

func (st *CommonStack[T]) Push(elem T) error {
	if st == nil {
		return stack.ErrorNilPointer
	}

	for {
		top := st.top.Load()
		newTop := newNode(elem, top)
		if st.top.CompareAndSwap(top, newTop) {
			return nil
		}
	}
}

func (st *CommonStack[T]) Pop() (T, error) {
	var elem T
	if st == nil {
		return elem, stack.ErrorNilPointer
	}
	for {
		top := st.top.Load()
		if top == nil {
			return elem, stack.ErrorEmptyStack
		} else {
			if st.top.CompareAndSwap(top, top.next.Load()) {
				return top.val, nil
			}
		}
	}
}

func (st *CommonStack[T]) Top() (T, error) {
	var elem T
	if st == nil {
		return elem, stack.ErrorNilPointer
	} else {
		top := st.top.Load()
		if top == nil {
			return elem, stack.ErrorEmptyStack
		} else {
			return top.val, nil
		}
	}
}

func (st *CommonStack[T]) Stringify() string {
	str := "> "
	if st == nil {
		return str
	}

	n := st.top.Load()
	for n != nil {
		str += fmt.Sprintf("%v ", n.val)
		n = n.next.Load()
	}
	return str
}