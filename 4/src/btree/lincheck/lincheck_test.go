//go:build !race
// +build !race

// The test takes too long under the race detector.

package lincheck

import (
	"fmt"
	"math/rand"
	"parallel-prog/4/btree"
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

	gen := btree.NewNoLockTree[int, int]()
	emptyGen := func() tree {
		return btree.NewNoLockTree[int, int]()
	}

	checker := NewSimpleTree[int, int]()
	emptyCheck := func() tree {
		return NewSimpleTree[int, int]()
	}
	for i := 0; i < diffSeedsNum; i++ {
		c := MakeChecker(gen, emptyGen, checker, emptyCheck, threadsNum, threadLen, int64(rg.Int()), 5*time.Second)
		trace := c.RunCheck(repeatRunNum)
		if trace != nil {
			fmt.Print(trace)
			return
		}
	}
	t.Errorf("Lincheck don't find errors in bad stack realization")

}
