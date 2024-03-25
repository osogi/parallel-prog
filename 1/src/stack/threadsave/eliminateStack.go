package threadsave

import (
	"errors"
	"fmt"
	"math/rand"
	"parallel-prog/1/stack"
	"sync/atomic"
	"time"
)

const spinMax = time.Millisecond * 512
const spinMin = time.Millisecond * 1
const spinDefault = time.Millisecond * 64

const spinCounterMax = 8
const factorCounterMax = 8

const EMPTY_COLLISION = uint32(0xffffffff)

var ErrorOpNoPerform = errors.New("operation wasn't performed")

type operation uint

const (
	PUSH operation = iota
	POP
	TOP
)

type threadInfo[T any] struct {
	id   uint
	op   operation
	elem T
}

type EliminateStackManager[T any] struct {
	top        atomic.Pointer[node[T]]
	collisions []atomic.Uint32
	threads    []atomic.Pointer[threadInfo[T]]
	maxThreads uint
	numThreads uint
}

type EliminateStack[T any] struct {
	manager  *EliminateStackManager[T]
	threadId uint

	spin        time.Duration
	spinCounter uint

	factor        float32
	factorCounter uint
}

func NewEliminateStackManager[T any](maxThreads uint) *EliminateStackManager[T] {
	esm := &EliminateStackManager[T]{collisions: make([]atomic.Uint32, maxThreads),
		threads:    make([]atomic.Pointer[threadInfo[T]], maxThreads),
		maxThreads: maxThreads,
		numThreads: 0,
		top:        atomic.Pointer[node[T]]{}}
	for i := range esm.collisions {
		esm.collisions[i].Store(EMPTY_COLLISION)
	}
	return esm
}

func (esm *EliminateStackManager[T]) GetStackForNewThread() (*EliminateStack[T], error) {
	es := &EliminateStack[T]{spin: spinDefault, spinCounter: spinCounterMax / 2,
		factor: 0.5, factorCounter: factorCounterMax / 2}
	if esm.numThreads < esm.maxThreads {
		es.manager = esm
		es.threadId = esm.numThreads
		esm.numThreads += 1
		return es, nil
	} else {
		return es, fmt.Errorf("can't add more threads for this stack")
	}
}

func (st *EliminateStackManager[T]) push(elem T) error {
	if st == nil {
		return stack.ErrorNilPointer
	}
	top := st.top.Load()
	newTop := newNode(elem, top)
	if st.top.CompareAndSwap(top, newTop) {
		return nil
	} else {
		return ErrorOpNoPerform
	}

}

func (st *EliminateStackManager[T]) pop() (T, error) {
	var elem T
	if st == nil {
		return elem, stack.ErrorNilPointer
	}
	top := st.top.Load()
	if top == nil {
		return elem, stack.ErrorEmptyStack
	} else {
		if st.top.CompareAndSwap(top, top.next.Load()) {
			return top.val, nil
		} else {
			return elem, ErrorOpNoPerform
		}
	}
}

func (st *EliminateStackManager[T]) getTop() (T, error) {
	var elem T
	if st == nil {
		return elem, stack.ErrorNilPointer
	} else {
		top := st.top.Load()
		if top == nil {
			return elem, stack.ErrorEmptyStack
		} else {
			return top.val, nil
		}
	}
}

func (st *EliminateStackManager[T]) stringify() string {
	str := "> "
	if st == nil {
		return str
	}

	n := st.top.Load()
	for n != nil {
		str += fmt.Sprintf("%v ", n.val)
		n = n.next.Load()
	}
	return str
}

func (st *EliminateStackManager[T]) tryPerformOp(inf *threadInfo[T]) (T, error) {
	var elem T
	switch inf.op {
	case PUSH:
		return elem, st.push(inf.elem)
	case POP:
		return st.pop()
	case TOP:
		return st.getTop()
	default:
		return elem, errors.ErrUnsupported
	}

}

func (loc_st *EliminateStack[T]) Push(elem T) error {
	if loc_st == nil {
		return stack.ErrorNilPointer
	}

	inf := &threadInfo[T]{id: loc_st.threadId, op: PUSH, elem: elem}
	_, err := loc_st.launchOp(inf)
	return err
}

func (loc_st *EliminateStack[T]) Pop() (T, error) {
	var elem T
	if loc_st == nil {
		return elem, stack.ErrorNilPointer
	}

	inf := &threadInfo[T]{id: loc_st.threadId, op: POP, elem: elem}
	elem, err := loc_st.launchOp(inf)
	return elem, err
}

func (loc_st *EliminateStack[T]) Top() (T, error) {
	var elem T
	if loc_st == nil {
		return elem, stack.ErrorNilPointer
	}

	inf := &threadInfo[T]{id: loc_st.threadId, op: TOP, elem: elem}
	elem, err := loc_st.launchOp(inf)
	return elem, err
}

func (loc_st *EliminateStack[T]) launchOp(inf *threadInfo[T]) (T, error) {
	var elem T
	if loc_st == nil {
		return elem, stack.ErrorNilPointer
	}

	st := loc_st.manager
	if st == nil {
		return elem, stack.ErrorNilPointer
	}

	for {
		// try just to perform op
		elem, err := st.tryPerformOp(inf)
		if !errors.Is(err, ErrorOpNoPerform) {
			return elem, err
		}

		// there started magic with elimination
		st.threads[inf.id].Store(inf)
		pos := loc_st.getPosition()
		him := st.collisions[pos].Load()
		for !st.collisions[pos].CompareAndSwap(him, uint32(loc_st.threadId)) {
			him = st.collisions[pos].Load()
		}

		// if found smth to collide with
		if him != EMPTY_COLLISION {
			loc_st.hitCollision()
			hisInf := st.threads[him].Load()
			//          | not already collided | reverse op
			if hisInf != nil && uint32(hisInf.id) == him && hisInf.op+inf.op == PUSH+POP {
				if !loc_st.isCollided(inf) {
					succ, elem := loc_st.tryCollision(inf, hisInf)
					if succ {
						return elem, nil
					} else {
						continue
					}
				} else {
					return loc_st.finishCollision(inf), nil
				}
			}
		} else {
			loc_st.skipCollision()
		}

		// try wait until smbdy will try collide with us
		time.Sleep(loc_st.spin)

		if loc_st.isCollided(inf) {
			loc_st.timeHitCollision()
			return loc_st.finishCollision(inf), nil
		} else {
			loc_st.timeSkipCollision()
		}

	}

}

// check is thread already collided, if not will remove it from pool
func (loc_st *EliminateStack[T]) isCollided(inf *threadInfo[T]) bool {
	st := loc_st.manager
	return !st.threads[loc_st.threadId].CompareAndSwap(inf, nil)
}

// need if was pop and push eliminate it
func (loc_st *EliminateStack[T]) finishCollision(inf *threadInfo[T]) T {
	var elem T
	st := loc_st.manager
	if inf.op == POP {
		elem = st.threads[inf.id].Load().elem
		st.threads[inf.id].Store(nil)
	}
	return elem
}

func (loc_st *EliminateStack[T]) tryCollision(selfInf *threadInfo[T], hisInf *threadInfo[T]) (bool, T) {
	var elem T
	st := loc_st.manager

	switch selfInf.op {
	case PUSH:
		if st.threads[hisInf.id].CompareAndSwap(hisInf, selfInf) {
			return true, elem
		} else { // already eliminated
			return false, elem
		}
	case POP:
		if st.threads[hisInf.id].CompareAndSwap(hisInf, nil) {
			elem = hisInf.elem
			st.threads[selfInf.id].Store(nil)
			return true, elem
		} else { // already eliminated
			return false, elem
		}
	default:
		return false, elem
	}

}

// functions for dynamic change of spin and factor value
func (loc_st *EliminateStack[T]) skipCollision() {
	loc_st.factorCounter += 1
	if loc_st.factorCounter >= factorCounterMax {
		loc_st.factorCounter = 1
		loc_st.factor /= 2
		loc_st.factor = max(0.0000001, loc_st.factor)
	}
}

func (loc_st *EliminateStack[T]) hitCollision() {
	loc_st.factorCounter -= 1
	if loc_st.factorCounter <= 0 {
		loc_st.factorCounter = factorCounterMax - 1
		loc_st.factor *= 2
		loc_st.factor = min(1, loc_st.factor)
	}
}

func (loc_st *EliminateStack[T]) timeSkipCollision() {
	loc_st.spinCounter += 1
	if loc_st.spinCounter >= spinCounterMax {
		loc_st.spinCounter = 1
		loc_st.spin *= 2
		loc_st.spin = min(spinMax, loc_st.spin)
	}
}

func (loc_st *EliminateStack[T]) timeHitCollision() {
	loc_st.spinCounter -= 1
	if loc_st.spinCounter <= 0 {
		loc_st.spinCounter = spinCounterMax - 1
		loc_st.spin /= 2
		loc_st.spin = max(spinMin, loc_st.spin)
	}
}

func (loc_st *EliminateStack[T]) getPosition() uint {
	var res uint
	m := int(loc_st.factor * float32(len(loc_st.manager.collisions)))
	m = min(m, len(loc_st.manager.collisions))
	if m != 0 {
		res = uint(rand.Intn(m))
	} else {
		res = 0
	}
	return res
}

func (loc_st *EliminateStack[T]) Stringify() string {
	if loc_st == nil {
		return ""
	}
	return loc_st.manager.stringify()
}
