package btree

import (
	"fmt"
	"sync"

	"golang.org/x/exp/constraints"
)

type node[K constraints.Ordered, V any] struct {
	Key         K
	Value       V
	Right, Left *node[K, V]
}

func newNode[K constraints.Ordered, V any](key K, value V) *node[K, V] {
	return &node[K, V]{Key: key, Value: value, Right: nil, Left: nil}
}

func (n *node[K, V]) find(key K) (*node[K, V], *node[K, V]) {
	var resParent *node[K, V] = nil

	res := n

	finded := false
	for !finded {
		if res == nil {
			break
		}

		if res.Key == key {
			finded = true
		} else {
			resParent = res
			if key < res.Key {
				res = res.Left
			} else {
				res = res.Right
			}
		}
	}
	return res, resParent
}

func (p *node[K, V]) insert(n *node[K, V]) error {
	errorNotNil := fmt.Errorf("can't insert node instead not nil node")
	if n.Key == p.Key {
		return fmt.Errorf("can't insert node to parent with same key")
	} else if n.Key < p.Key {
		if p.Left != nil {
			return errorNotNil
		} else {
			p.Left = n
		}
	} else {
		if p.Right != nil {
			return errorNotNil
		} else {
			p.Right = n
		}
	}

	return nil
}

type HardLockTree[K constraints.Ordered, V any] struct {
	root  *node[K, V]
	mutex sync.Mutex
}

func NewHardLockTree[K constraints.Ordered, V any]() *HardLockTree[K, V] {
	return &HardLockTree[K, V]{root: nil, mutex: sync.Mutex{}}
}

func (t *HardLockTree[K, V]) _find(key K) (*node[K, V], *node[K, V]) {
	if t.root == nil {
		return nil, nil
	} else {
		return t.root.find(key)
	}
}

func (t *HardLockTree[K, V]) Find(key K) (V, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	n, _ := t._find(key)

	if n == nil {
		var emptyValue V
		return emptyValue, ErrorNodeNotFound
	} else {
		return n.Value, nil
	}

}

func (t *HardLockTree[K, V]) Insert(key K, value V) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	new := newNode(key, value)
	if t.root == nil {
		t.root = new
	} else {
		n, p := t.root.find(key)
		if n != nil {
			return ErrorSameKey
		}
		return p.insert(new)
	}
	return nil
}

func (t *HardLockTree[K, V]) replaceNode(parent *node[K, V], prevNode *node[K, V], newNode *node[K, V]) {
	if parent == nil {
		t.root = newNode
	} else {
		if parent.Left == prevNode {
			parent.Left = newNode
		} else {
			parent.Right = newNode
		}
	}
}

func (t *HardLockTree[K, V]) _delete(target *node[K, V], parent *node[K, V]) error {
	if target == nil {
		return ErrorNodeNotFound
	} else {
		if target.Left == nil {
			t.replaceNode(parent, target, target.Right)
		} else if target.Right == nil {
			t.replaceNode(parent, target, target.Left)
		} else {
			prev := target
			cur := target.Left
			for cur.Right != nil {
				prev = cur
				cur = cur.Right
			}

			err := t._delete(cur, prev)
			if err != nil {
				return err
			}
			cur.Right = target.Right
			t.replaceNode(parent, target, cur)
		}
	}

	return nil
}

func (t *HardLockTree[K, V]) Delete(key K) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	n, p := t._find(key)

	if n == nil {
		return ErrorNodeNotFound
	} else {
		return t._delete(n, p)
	}
}

// func (t *HardLockTree[K, V]) Print() {
// 	cur := make(chan *node[K, V], 300)
// 	next := make(chan *node[K, V], 300)

// 	cur <- t.root
// 	nullOnly := false

// 	i := 0
// 	for !nullOnly {
// 		nullOnly = true
// 		for j := 0; j < 1<<i; j++ {
// 			n := <-cur
// 			fmt.Printf("%v ", n)
// 			if n != nil {
// 				nullOnly = false
// 				next <- n.Left
// 				next <- n.Right
// 			} else {
// 				next <- nil
// 				next <- nil
// 			}
// 		}

// 		cur = next
// 		next = make(chan *node[K, V], 300)
// 		i++
// 		fmt.Printf("\n")
// 	}
// }
