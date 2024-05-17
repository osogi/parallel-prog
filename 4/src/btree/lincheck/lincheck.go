package lincheck

import (
	"fmt"
	"math/rand"
	"parallel-prog/4/btree"
	"sync"
	"time"
)

const MAX_INT = 16

type insertEvent struct {
	kv  int
	res error

	resCheck error
}

func (e *insertEvent) generate(rg *rand.Rand) {
	e.kv = rg.Intn(MAX_INT)
}
func (e *insertEvent) execute(t btree.Btree[int, int]) {
	e.res = t.Insert(e.kv, e.kv)
}
func (e *insertEvent) check(t btree.Btree[int, int]) bool {
	e.resCheck = t.Insert(e.kv, e.kv)
	return e.res == e.resCheck
}
func (e *insertEvent) rollback(t btree.Btree[int, int]) {
	if e.resCheck == nil {
		t.Delete(e.kv)
	}
}
func (e *insertEvent) stringify() string {
	return fmt.Sprintf("Insert(%d): | err: %v", e.kv, e.res)
}
func (e *insertEvent) clean() {
	*e = insertEvent{kv: e.kv}
}

type deleteEvent struct {
	kv  int
	res error

	resCheck error
}

func (e *deleteEvent) generate(rg *rand.Rand) {
	e.kv = rg.Intn(MAX_INT)
}
func (e *deleteEvent) execute(t btree.Btree[int, int]) {
	e.res = t.Delete(e.kv)
}
func (e *deleteEvent) check(t btree.Btree[int, int]) bool {
	e.resCheck = t.Delete(e.kv)
	return e.res == e.resCheck
}
func (e *deleteEvent) rollback(t btree.Btree[int, int]) {
	if e.resCheck == nil {
		t.Insert(e.kv, e.kv)
	}
}
func (e *deleteEvent) stringify() string {
	return fmt.Sprintf("Delete(%d): | err: %v", e.kv, e.res)
}
func (e *deleteEvent) clean() {
	*e = deleteEvent{kv: e.kv}
}

type findEvent struct {
	kv  int
	res error
}

func (e *findEvent) generate(rg *rand.Rand) {
	e.kv = rg.Intn(MAX_INT)
}
func (e *findEvent) execute(t btree.Btree[int, int]) {
	_, e.res = t.Find(e.kv)
}
func (e *findEvent) check(t btree.Btree[int, int]) bool {
	cp := *e
	cp.execute(t)
	return cp == *e
}
func (e *findEvent) rollback(t btree.Btree[int, int]) {
}
func (e *findEvent) stringify() string {
	return fmt.Sprintf("Find(%d): | err: %v", e.kv, e.res)
}
func (e *findEvent) clean() {
	*e = findEvent{kv: e.kv}
}

type event interface {
	generate(*rand.Rand)
	execute(btree.Btree[int, int])
	check(btree.Btree[int, int]) bool
	stringify() string
	clean()

	// this is a little bit strange too. Maybe we should use some other method
	rollback(btree.Btree[int, int])
}

type thread struct {
	events  []event
	curstep int
}

func newThread(size int, rg *rand.Rand) *thread {
	events := make([]event, size)
	for i := 0; i < size; i++ {
		var e event
		typeint := rg.Intn(3)
		switch typeint {
		case 0:
			e = &insertEvent{}
		case 1:
			e = &deleteEvent{}
		case 2:
			e = &findEvent{}
		}
		e.generate(rg)
		events[i] = e
	}
	return &thread{events: events, curstep: 0}
}

func (th *thread) execute(t btree.Btree[int, int]) {

	for i, ev := range th.events {
		th.curstep = i
		ev.execute(t)
		sl := rand.Intn(64)
		if sl < 10 {
			time.Sleep(time.Duration(sl) * time.Millisecond)
		}
	}
	th.curstep = len(th.events)
}

func (th *thread) clean() {
	for _, ev := range th.events {
		ev.clean()
	}
}

type Checker struct {
	threads           []*thread
	generateTree      btree.Btree[int, int]
	checkTree         btree.Btree[int, int]
	timeout           time.Duration
	emptyCheckTree    func() btree.Btree[int, int]
	emptyGenerateTree func() btree.Btree[int, int]
}

func MakeChecker(generateTree btree.Btree[int, int], emptyGenerateTree func() btree.Btree[int, int], checkTree btree.Btree[int, int], emptyCheckTree func() btree.Btree[int, int], threadNum int, threadSize int, randSeed int64, timeout time.Duration) Checker {
	randGen := rand.New(rand.NewSource(randSeed))
	threads := make([]*thread, threadNum)
	for i := 0; i < threadNum; i++ {
		threads[i] = newThread(threadSize, randGen)
	}
	return Checker{threads: threads, generateTree: generateTree, checkTree: checkTree, timeout: timeout, emptyCheckTree: emptyCheckTree, emptyGenerateTree: emptyGenerateTree}
}

func (c *Checker) run(timeout time.Duration) (bool, []int) {
	randevu := sync.WaitGroup{}
	randevu.Add(len(c.threads))
	isTimeEnd := make(chan bool, len(c.threads))

	for _, th := range c.threads {
		th.clean()
		t := th

		go func() {
			ended := make(chan bool, 1)
			randevu.Done()
			randevu.Wait()

			go func() {
				t.execute(c.generateTree)
				ended <- true
			}()
			select {
			case <-ended:
				isTimeEnd <- false
			case <-time.After(timeout):
				isTimeEnd <- true
			}
		}()
	}

	noTimeouts := true
	for i := 0; i < len(c.threads); i++ {
		isTimeout := <-isTimeEnd
		noTimeouts = noTimeouts && !isTimeout
	}

	lastSteps := make([]int, len(c.threads))
	for i, th := range c.threads {
		lastSteps[i] = th.curstep
	}

	return noTimeouts, lastSteps
}

func sum(a []int) int {
	n := 0
	for _, v := range a {
		n += v
	}
	return n
}

func (c *Checker) checkStep(thread_exe *[]int) (bool, []int) {
	end := 0
	// fmt.Printf("%v\n", thread_exe)
	find := false
	retArr := make([]int, len(*thread_exe))
	copy(retArr, *thread_exe)
	for i, num := range *thread_exe {
		if num == len(c.threads[i].events) {
			end += 1
		} else {
			if c.threads[i].events[num].check(c.checkTree) {
				(*thread_exe)[i] = num + 1
				v, arr := c.checkStep(thread_exe)
				(*thread_exe)[i] = num
				if !v {
					if sum(retArr) < sum(arr) {
						retArr = arr
					}
				} else {
					find = true
				}
			}
			c.threads[i].events[num].rollback(c.checkTree)
		}
		if find {
			break
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
	return c.checkStep(&thread_exe)
}

func (c *Checker) RunCheck(num int) error {
	var err error = nil
	var passCheck, passTimeout bool
	var failSteps []int
	for i := 1; i <= num; i++ {
		c.generateTree = c.emptyGenerateTree()
		passTimeout, failSteps = c.run(c.timeout)
		if passTimeout {
			c.checkTree = c.emptyCheckTree()
			passCheck, failSteps = c.check()
		}

		// fmt.Print("\n\n\n")
		if !passTimeout || !passCheck {
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
			var errMsg string
			switch {
			case !passTimeout:
				errMsg = "timeout was reached. Trace of execution:"
			case !passCheck:
				errMsg = "caught a case of impossible linearization in this trace"
			default:
				errMsg = "osogi is bad coder"
			}
			err = fmt.Errorf("[%d/%d]: %s\n%s", i, num, errMsg, trace)

			break
		}
	}
	return err
}
