//go:build !race
// +build !race

// The test takes too long under the race detector.

package lincheck_test

import (
	"fmt"
	"math/rand"
	"parallel-prog/4/btree"
	"parallel-prog/4/btree/lincheck"
	"parallel-prog/4/btree/trees"
	"testing"
	"time"
)

const seed = 0xaaaa
const diffSeedsNum = 100
const repeatRunNum = 10
const threadsNum = 5
const threadLen = 3

func TestLincheck(t *testing.T) {
	rg := rand.New(rand.NewSource(seed))

	gen := trees.NewNoLockTree[int, int]()
	emptyGen := func() btree.Btree[int, int] {
		return trees.NewNoLockTree[int, int]()
	}

	checker := lincheck.NewSimpleTree[int, int]()
	emptyCheck := func() btree.Btree[int, int] {
		return lincheck.NewSimpleTree[int, int]()
	}
	for i := 0; i < diffSeedsNum; i++ {
		c := lincheck.MakeChecker(gen, emptyGen, checker, emptyCheck, threadsNum, threadLen, int64(rg.Int()), 5*time.Second)
		trace := c.RunCheck(repeatRunNum)
		if trace != nil {
			fmt.Print(trace)
			return
		}
	}
	t.Errorf("Lincheck don't find errors in bad stack realization")

}
