package lincheck

import (
	"fmt"
	"math/rand"
	"parallel-prog/1/stack"
	"parallel-prog/1/stack/unsafe"
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
	bad := unsafe.NewUnsafeStack[int]()
	rg := rand.New(rand.NewSource(seed))

	sts := make([]stack.Stack[int], threadsNum)
	for a := 0; a < threadsNum; a++ {
		sts[a] = bad
	}
	for i := 0; i < diffSeedsNum; i++ {
		c := MakeChecker(sts, threadLen, int64(rg.Int()), 5*time.Second)
		trace := c.RunCheck(repeatRunNum)
		if trace != nil {
			fmt.Print(trace)
			return
		}
	}
	t.Errorf("Lincheck don't find errors in bad stack realization")

}
