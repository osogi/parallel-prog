package lincheck

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"parallel-prog/1/stack"
	"sync"
	"time"
)

type pushEvent struct {
	elem int
	res  error
}

func (pe *pushEvent) generate(rg *rand.Rand) {
	pe.elem = rg.Intn(100)
}

func (pe *pushEvent) execute(st stack.Stack[int]) {
	pe.res = st.Push(pe.elem)
}
func (pe *pushEvent) check(st stack.Stack[int]) bool {
	cp := *pe
	cp.execute(st)
	return cp == *pe
}
func (pe *pushEvent) rollback(st stack.Stack[int]) {
	st.Pop()
}
func (pe *pushEvent) stringify() string {
	return fmt.Sprintf("Push(%d): | err: %v", pe.elem, pe.res)
}

type popEvent struct {
	elem    int
	err     error
	checked *popEvent
}

func (pe *popEvent) generate(rg *rand.Rand) {}
func (pe *popEvent) execute(st stack.Stack[int]) {
	pe.elem, pe.err = st.Pop()
}
func (pe *popEvent) check(st stack.Stack[int]) bool {
	cp := *pe
	cp.execute(st)

	// this is realy annoing and should be changed
	pe.checked = &cp

	return ((cp.err == pe.err) && (cp.elem == pe.elem))
}
func (pe *popEvent) rollback(st stack.Stack[int]) {
	if pe.checked.err == nil {
		st.Push(pe.checked.elem)
	}
}
func (pe *popEvent) stringify() string {
	return fmt.Sprintf("Pop(): %d | err: %v", pe.elem, pe.err)
}

type topEvent struct {
	elem int
	err  error
}

func (te *topEvent) generate(rg *rand.Rand) {}
func (te *topEvent) execute(st stack.Stack[int]) {
	te.elem, te.err = st.Top()
}
func (te *topEvent) check(st stack.Stack[int]) bool {
	cp := *te
	cp.execute(st)
	return cp == *te
}
func (te *topEvent) rollback(st stack.Stack[int]) {}
func (te *topEvent) stringify() string {
	return fmt.Sprintf("Top(): %d | err: %v", te.elem, te.err)
}

type event interface {
	generate(*rand.Rand)
	execute(stack.Stack[int])
	check(stack.Stack[int]) bool
	stringify() string

	// this is a little bit strange too. Maybe we should use some other method
	rollback(stack.Stack[int])
}

type thread struct {
	events []event
}

func makeThread(size int, rg *rand.Rand) thread {
	events := make([]event, size)
	for i := 0; i < size; i++ {
		var e event
		typeint := rg.Intn(3)
		switch typeint {
		case 0:
			e = &pushEvent{}
		case 1:
			e = &popEvent{}
		case 2:
			e = &topEvent{}
		}
		e.generate(rg)
		events[i] = e
	}
	return thread{events: events}
}

func (th *thread) execute(st stack.Stack[int]) {

	for _, ev := range th.events {
		ev.execute(st)
		sl := rand.Intn(512)
		if sl < 100 {
			time.Sleep(time.Duration(sl) * time.Millisecond)
		}
	}
}

type Checker struct {
	threads []thread
	st      stack.Stack[int]
}

func MakeChecker(st stack.Stack[int], threadNum int, threadSize int, randSeed int64) Checker {
	randGen := rand.New(rand.NewSource(randSeed))
	threads := make([]thread, threadNum)
	for i := 0; i < threadNum; i++ {
		threads[i] = makeThread(threadSize, randGen)
	}
	return Checker{threads: threads, st: st}
}

func (c *Checker) run() {
	wg := sync.WaitGroup{}
	wg.Add(len(c.threads))
	for _, th := range c.threads {
		go func() {
			th.execute(c.st)
			wg.Done()
		}()
	}
	wg.Wait()
}

func sum(a []int) int {
	n := 0
	for _, v := range a {
		n += v
	}
	return n
}

func (c *Checker) checkStep(thread_exe []int) (bool, []int) {
	var find = false
	retArr := thread_exe
	end := 0
	//fmt.Printf("%v %s\n", thread_exe, c.st.Stringify())

	for i, num := range thread_exe {
		if num == len(c.threads[i].events) {
			end += 1
		} else {
			if c.threads[i].events[num].check(c.st) {
				new_thread_exe := make([]int, len(thread_exe))
				copy(new_thread_exe, thread_exe)
				new_thread_exe[i] = num + 1
				v, arr := c.checkStep(new_thread_exe)
				if !v {
					if sum(retArr) < sum(arr) {
						retArr = arr
					}
				} else {
					find = v
				}
			}
			c.threads[i].events[num].rollback(c.st)
		}
	}
	if end == len(c.threads) {
		find = true
	}
	return find, retArr
}

func (c *Checker) check() (bool, []int) {
	thread_exe := make([]int, len(c.threads))
	for i := range thread_exe {
		thread_exe[i] = 0
	}

	return c.checkStep(thread_exe)
}

func emptyStack[T any](st stack.Stack[T]) {
	var err error
	err = nil
	for err == nil {
		_, err = st.Pop()
	}

	if errors.Is(err, stack.ErrorEmptyStack) {
		return
	} else {
		log.Fatal(err)
	}
}

func (c *Checker) RunCheck(num int) error {
	var err error = nil
	for i := 1; i <= num; i++ {
		emptyStack(c.st)
		c.run()
		emptyStack(c.st)
		pass, failSteps := c.check()
		// fmt.Print("\n\n\n")
		if !pass {
			trace := ""
			for i, th := range c.threads {
				trace += fmt.Sprintf("Thread %d:\n", i)
				for evNum, ev := range th.events {
					pointer := "  "
					if failSteps[i] == evNum {
						pointer = "->"
					}
					trace += fmt.Sprintf("\t%s%s\n", pointer, ev.stringify())
				}
			}

			err = fmt.Errorf("[%d/%d]: caught a case of impossible linearization in this trace\n%s", i, num, trace)
			break
		}
	}
	return err
}
