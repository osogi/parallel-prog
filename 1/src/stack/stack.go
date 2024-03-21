package stack

type Stack[T any] interface {
	Push(T) error
	Pop() (T, error)
	Top() (T, error)
	Stringify() string
}
