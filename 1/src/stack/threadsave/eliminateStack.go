package threadsave

import (
	"sync/atomic"
)

type EliminateStack[T any] struct {
	top atomic.Pointer[node[T]]
}
