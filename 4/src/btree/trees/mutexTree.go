package trees

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type mutexTree[K constraints.Ordered, V any] struct {
	tree  *CpsTree[K, V, *mutexNode[K, V]]
	mutex sync.Mutex
}

func newMutexTree[K constraints.Ordered, V any]() *mutexTree[K, V] {
	return &mutexTree[K, V]{tree: NewCpsTree[K, V, *mutexNode[K, V]](nil), mutex: sync.Mutex{}}
}

func (t *mutexTree[K, V]) getNodeMutex(cur *mutexNode[K, V], asParent bool) *sync.Mutex {
	if cur.IsNil() {
		if asParent {
			return &t.mutex
		} else {
			return nil
		}
	} else {
		return &(cur.mutex)
	}
}

func (t *mutexTree[K, V]) lockNode(cur *mutexNode[K, V], asParent bool) {
	// str := ""
	// if cur == nil {
	// 	if asParent {
	// 		str = "tree"
	// 	} else {
	// 		str = "<nil>"
	// 	}
	// } else {
	// 	str = fmt.Sprintf("[%v]", cur.GetKey())
	// }
	// fmt.Printf("[+]: %v\n", str)
	mut := t.getNodeMutex(cur, asParent)
	if mut != nil {
		mut.Lock()
	}
}

func (t *mutexTree[K, V]) unlockNode(cur *mutexNode[K, V], asParent bool) {
	// str := ""
	// if cur == nil {
	// 	if asParent {
	// 		str = "tree"
	// 	} else {
	// 		str = "<nil>"
	// 	}
	// } else {
	// 	str = fmt.Sprintf("[%v]", cur.GetKey())
	// }
	// fmt.Printf("[ ]: %v\n", str)
	mut := t.getNodeMutex(cur, asParent)
	if mut != nil {
		mut.Unlock()
	}
}
