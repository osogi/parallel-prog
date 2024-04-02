package main

import (
	"fmt"
	"math/rand"
	"parallel-prog/1/stack"
	"parallel-prog/1/stack/threadsave"
	"sync"
	"testing"
	"time"
)

func stackWorkload(st stack.Stack[int], opsCount int, pushChance float64, otherWorkload time.Duration) {
	for i := 0; i < opsCount; i++ {
		rnd := rand.Float64()
		switch {
		case rnd < pushChance:
			st.Push(rand.Intn(100))
		default:
			st.Pop()
		}
		time.Sleep(otherWorkload)
		// end := time.Now().Add(otherWorkload)
		// c := 1
		// for !end.Before(time.Now()) {
		// 	c = c + i*i
		// }

	}
}

func runStackBenchmark(b *testing.B, sts []stack.Stack[int], opsCount int, pushChance float64, otherWorkload time.Duration) {
	threadsNum := len(sts)

	start := sync.WaitGroup{}
	end := sync.WaitGroup{}
	start.Add(threadsNum + 1)
	end.Add(threadsNum)

	for _, st := range sts {
		go func() {
			start.Done()
			start.Wait()

			stackWorkload(st, opsCount, pushChance, otherWorkload)

			end.Done()
		}()
	}
	b.StartTimer()
	start.Done()
	start.Wait()

	end.Wait()

}

func runEliminateStackBench(b *testing.B, threadsNum int, opsCount int, pushChance float64, otherWorkload time.Duration) {

	b.Run(fmt.Sprintf("Eliminate(goroutines=%d,ops=%d,pushChance=%f,sleep=%v)", threadsNum, opsCount, pushChance, otherWorkload),
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				stm := threadsave.NewEliminateStackManager[int](uint(threadsNum))

				sts := make([]stack.Stack[int], threadsNum)
				for a := 0; a < threadsNum; a++ {
					var err error
					sts[a], err = stm.GetStackForNewThread()
					if err != nil {
						b.Error(err.Error())
					}
				}
				runStackBenchmark(b, sts, opsCount, pushChance, otherWorkload)
			}
		})
}

func runCommonStackBench(b *testing.B, threadsNum int, opsCount int, pushChance float64, otherWorkload time.Duration) {
	b.Run(fmt.Sprintf("Common(goroutines=%d,ops=%d,pushChance=%f,sleep=%v)", threadsNum, opsCount, pushChance, otherWorkload),
		func(b *testing.B) {
			b.StopTimer()
			for i := 0; i < b.N; i++ {
				st := threadsave.NewCommonStack[int]()

				sts := make([]stack.Stack[int], threadsNum)
				for a := 0; a < threadsNum; a++ {
					sts[a] = st
				}

				runStackBenchmark(b, sts, opsCount, pushChance, otherWorkload)
				//fmt.Printf("%v/%v\n", i, b.N)

			}
		})
}

func Benchmark(b *testing.B) {

	runCommonStackBench(b, 10000, 1000, 0.5, 0*time.Microsecond)
	runEliminateStackBench(b, 10000, 1000, 0.5, 0*time.Microsecond)

	runCommonStackBench(b, 4000, 1000, 0.5, 5*time.Microsecond)
	runEliminateStackBench(b, 4000, 1000, 0.5, 5*time.Microsecond)

	runCommonStackBench(b, 2000, 1000, 0.5, 50*time.Microsecond)
	runEliminateStackBench(b, 2000, 1000, 0.5, 50*time.Microsecond)

	runCommonStackBench(b, 2000, 1000, 0.5, 500*time.Microsecond)
	runEliminateStackBench(b, 2000, 1000, 0.5, 500*time.Microsecond)

	runCommonStackBench(b, 10000, 5000, 0.25, 0*time.Microsecond)
	runEliminateStackBench(b, 10000, 5000, 0.25, 0*time.Microsecond)

}
