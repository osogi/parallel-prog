//go:build !race
// +build !race

package threadsave

import (
	"math/rand"
	"parallel-prog/1/stack"
	"parallel-prog/1/stack/lincheck"
	"testing"
	"time"
)

const seed = 0xaaaa
const diffSeedsNum = 100
const repeatRunNum = 10
const threadsNum = 6
const threadLen = 3

// The test takes too long under the race detector.
func TestLincheck(t *testing.T) {
	st := NewCommonStack[int]()
	rg := rand.New(rand.NewSource(seed))

	sts := make([]stack.Stack[int], threadsNum)
	for a := 0; a < threadsNum; a++ {
		sts[a] = st
	}
	for i := 0; i < diffSeedsNum; i++ {
		c := lincheck.MakeChecker(sts, threadLen, int64(rg.Int()), 5*time.Second)
		err := c.RunCheck(repeatRunNum)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	}

}
