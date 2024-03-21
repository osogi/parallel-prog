package stack

import "fmt"

type unsafeNode[T any] struct {
	val  T
	next *unsafeNode[T]
}

type UnsafeStack[T any] struct {
	top *unsafeNode[T]
}

func newUnsafeNode[T any](val T, ptr *unsafeNode[T]) *unsafeNode[T] {
	return &unsafeNode[T]{val, ptr}
}

func NewUnsafeStack[T any]() *UnsafeStack[T] {
	return &UnsafeStack[T]{nil}
}

func (st *UnsafeStack[T]) Push(elem T) error {
	if st == nil {
		return ErrorNilPointer
	}

	if st.top == nil {
		st.top = newUnsafeNode(elem, nil)
	} else {
		buf := newUnsafeNode(elem, st.top)
		st.top = buf
	}
	return nil
}

func (st *UnsafeStack[T]) Pop() (T, error) {
	var elem T
	if st == nil {
		return elem, ErrorNilPointer
	}

	if st.top == nil {
		return elem, ErrorEmptyStack
	} else {
		elem = st.top.val
		st.top = st.top.next
		return elem, nil
	}
}

func (st *UnsafeStack[T]) Top() (T, error) {
	var elem T
	if st == nil {
		return elem, ErrorNilPointer
	} else {
		top := st.top
		if top == nil {
			return elem, ErrorEmptyStack
		} else {
			return top.val, nil
		}
	}
}

func (st *UnsafeStack[T]) Stringify() string {
	str := "> "
	n := st.top
	for n != nil {
		str += fmt.Sprintf("%v ", n.val)
		n = n.next
	}
	return str
}
