package threadsave

import (
	"math/rand"
	"parallel-prog/1/stack"
	"testing"
)

const opsCount = 10000
const goroutinesCount = 1000

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

func TestRaceCommon(t *testing.T) {
	st := NewCommonStack[int]()
	for i := 0; i < goroutinesCount; i++ {
		go stackWorkload(st, opsCount)
	}

}

func TestRaceEliminate(t *testing.T) {
	stm := NewEliminateStackManager[int](goroutinesCount)

	for i := 0; i < goroutinesCount; i++ {
		st, err := stm.GetStackForNewThread()
		if err != nil {
			t.Error(err.Error())
		}
		go stackWorkload(st, opsCount)
	}

}
