package lincheck

import (
	"fmt"
	"math/rand"
	"parallel-prog/1/stack"
	"testing"
)

func TestLincheck(t *testing.T) {
	bad := stack.NewUnsafeStack[int]()
	rg := rand.New(rand.NewSource(0xaaaa))
	for i := 0; i < 100; i++ {
		c := MakeChecker(bad, 2, 10, int64(rg.Int()))
		trace := c.RunCheck(10)
		if trace != nil {
			fmt.Print(trace)
			return
		}
	}
	t.Errorf("Lincheck don't find errors in bad stack realization")
}
