package threadsave

import (
	"math/rand"
	"parallel-prog/1/stack"
	"testing"
)

func stackWorkload(st stack.Stack[int], opsCount int) {
	for i := 0; i < opsCount; i++ {
		switch rand.Intn(3) {
		case 0:
			st.Push(rand.Intn(100))
		case 1:
			st.Pop()
		case 2:
			st.Top()
		}
	}
}

func TestRace(t *testing.T) {
	opsCount := 10000
	goroutinesCount := 1000
	st := NewCommonStack[int]()
	for i := 0; i < goroutinesCount; i++ {
		go stackWorkload(st, opsCount)
	}

}
