package btree

import (
	"fmt"
	"sync"

	"golang.org/x/exp/constraints"
)

type HardLockTree[K constraints.Ordered, V any] struct {
	tree  *CpsTree[K, V, *commonNode[K, V]]
	mutex sync.Mutex
}

func NewHardLockTree[K constraints.Ordered, V any]() *HardLockTree[K, V] {
	return &HardLockTree[K, V]{tree: NewCpsTree[K, V, *commonNode[K, V]](nil), mutex: sync.Mutex{}}
}

func emptyFunc[K constraints.Ordered, V any](a *commonNode[K, V], b *commonNode[K, V], c *commonNode[K, V]) {
}

func (t *HardLockTree[K, V]) subFind(cur *commonNode[K, V], parent *commonNode[K, V], key K) (*commonNode[K, V], *commonNode[K, V]) {
	return t.tree.NodeFind(cur, parent, nil, key, emptyFunc)
}
func (t *HardLockTree[K, V]) subInsert(cur *commonNode[K, V], parent *commonNode[K, V], newNode *commonNode[K, V]) {
	t.tree.NodeInsert(cur, parent, newNode, t.subFind)
}

func (t *HardLockTree[K, V]) Find(key K) (V, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var nilValue V

	n, _ := t.subFind(t.tree.root, nil, key)
	if n.IsNil() {
		return nilValue, ErrorNodeNotFound
	} else {
		return n.GetValue(), nil
	}
}

func (t *HardLockTree[K, V]) Insert(key K, value V) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	newNode := newCommonNode(key, value)
	_, _, err := t.tree.NodeInsert(t.tree.root, nil, newNode, t.subFind)
	return err
}

func (t *HardLockTree[K, V]) Delete(key K) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	_, _, err := t.tree.NodeDelete(t.tree.root, nil, key, t.subFind, t.subInsert)
	return err
}

func (t *HardLockTree[K, V]) Print() {
	cur := make(chan *commonNode[K, V], 300)
	next := make(chan *commonNode[K, V], 300)

	cur <- t.tree.root
	nullOnly := false

	i := 0
	for !nullOnly {
		nullOnly = true
		for j := 0; j < 1<<i; j++ {
			n := <-cur
			fmt.Printf("%v ", n)
			if n != nil {
				nullOnly = false
				next <- n.Left
				next <- n.Right
			} else {
				next <- nil
				next <- nil
			}
		}

		cur = next
		next = make(chan *commonNode[K, V], 300)
		i++
		fmt.Printf("\n")
	}
}
