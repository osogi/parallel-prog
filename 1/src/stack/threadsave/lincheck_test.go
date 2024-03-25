//go:build !race
// +build !race

// The tests takes too long under the race detector.

package threadsave

import (
	"math/rand"
	"parallel-prog/1/stack"
	"parallel-prog/1/stack/lincheck"
	"testing"
	"time"
)

func RunLincheckTest(t *testing.T, sts []stack.Stack[int], threadLen int, seed int64, diffSeedsNum int, repeatRunNum int) {
	rg := rand.New(rand.NewSource(seed))
	for i := 0; i < diffSeedsNum; i++ {
		c := lincheck.MakeChecker(sts, threadLen, int64(rg.Int()), 5*time.Second)
		err := c.RunCheck(repeatRunNum)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	}
}

func RunLincheckTestCommon(t *testing.T, threadsNum int, threadLen int, seed int64, diffSeedsNum int, repeatRunNum int) {
	st := NewCommonStack[int]()

	sts := make([]stack.Stack[int], threadsNum)
	for a := 0; a < threadsNum; a++ {
		sts[a] = st
	}
	RunLincheckTest(t, sts, threadLen, seed, diffSeedsNum, repeatRunNum)
}

func RunLincheckTestEliminate(t *testing.T, threadsNum int, threadLen int, seed int64, diffSeedsNum int, repeatRunNum int) {
	stm := NewEliminateStackManager[int](uint(threadsNum))

	sts := make([]stack.Stack[int], threadsNum)
	for a := 0; a < threadsNum; a++ {
		var err error
		sts[a], err = stm.GetStackForNewThread()
		if err != nil {
			t.Error(err.Error())
		}
	}
	RunLincheckTest(t, sts, threadLen, seed, diffSeedsNum, repeatRunNum)
}

func TestLincheckCommon_6threads_3ops(t *testing.T) {
	RunLincheckTestCommon(t, 6, 3, 0xaaaa, 100, 10)
}

func TestLincheckEliminate_6threads_3ops(t *testing.T) {
	RunLincheckTestEliminate(t, 6, 3, 0xaaaa, 100, 10)
}

func TestLincheckEliminate_3threads_10ops(t *testing.T) {
	RunLincheckTestEliminate(t, 3, 10, 0xaaaa, 100, 10)
}

func TestLincheckEliminate_2threads_20ops(t *testing.T) {
	RunLincheckTestEliminate(t, 2, 20, 0xaaaa, 100, 10)
}
