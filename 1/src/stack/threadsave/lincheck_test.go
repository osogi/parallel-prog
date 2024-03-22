//go:build !race
// +build !race

package threadsave

import (
	"math/rand"
	"parallel-prog/1/stack/lincheck"
	"testing"
	"time"
)

// The test takes too long under the race detector.
func TestLincheck(t *testing.T) {
	st := NewCommonStack[int]()
	rg := rand.New(rand.NewSource(0xaaaa))
	for i := 0; i < 100; i++ {
		c := lincheck.MakeChecker(st, 6, 3, int64(rg.Int()), 5*time.Second)
		err := c.RunCheck(10)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	}

}
