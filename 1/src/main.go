package main

import (
	"fmt"
	"math/rand"
	"parallel-prog/1/stack"
	"parallel-prog/1/stack/lincheck"
	"time"
)

func main() {
	bad := stack.NewCommonStack[int]()
	rg := rand.New(rand.NewSource(0xdead))
	for i := 0; i < 50; i++ {
		c := lincheck.MakeChecker(bad, 3, 10, int64(rg.Int()), 5*time.Second)
		a := c.RunCheck(10)
		fmt.Printf("%v\n", a)
	}
}
