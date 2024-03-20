package stack

type node[T any] struct {
	val  T
	next *node[T]
}

type CommonStack[T any] struct {
	top *node[T]
}

func newNode[T any](val T, ptr *node[T]) *node[T] {
	return &node[T]{val, ptr}
}

func NewCommonStack[T any]() *CommonStack[T] {
	return &CommonStack[T]{nil}
}

func (st *CommonStack[T]) Push(elem T) error {
	if st == nil {
		return ErrorNilPointer
	}

	if st.top == nil {
		st.top = newNode(elem, nil)
	} else {
		buf := newNode(elem, st.top)
		st.top = buf
	}
	return nil
}

func (st *CommonStack[T]) Pop() (T, error) {
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

func (st *CommonStack[T]) Top() (T, error) {
	var elem T
	if st == nil {
		return elem, ErrorNilPointer
	} else {
		return st.top.val, nil
	}
}
