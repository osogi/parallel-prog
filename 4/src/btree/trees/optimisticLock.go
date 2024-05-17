package trees

import (
	"fmt"
	"log"
	"parallel-prog/4/btree"

	"golang.org/x/exp/constraints"
)

type OptimisticLockTree[K constraints.Ordered, V any] struct {
	mtree *mutexTree[K, V]
}

func NewOptimisticLockTree[K constraints.Ordered, V any]() *OptimisticLockTree[K, V] {
	return &OptimisticLockTree[K, V]{mtree: newMutexTree[K, V]()}
}

func (extenalTree *OptimisticLockTree[K, V]) subFind(cur *mutexNode[K, V], parent *mutexNode[K, V], key K) (*mutexNode[K, V], *mutexNode[K, V]) {
	t := extenalTree.mtree

	_, _ = cur, parent
	preStep := func(cur *mutexNode[K, V], parent *mutexNode[K, V], grandpar *mutexNode[K, V]) {
	}

	for {
		fcur, fparent := t.tree.NodeFind(t.tree.root, nil, nil, key, preStep)
		t.lockNode(fparent, true)
		t.lockNode(fcur, false)

		ncur, nparent := t.tree.NodeFind(t.tree.root, nil, nil, key, preStep)
		if fcur.IsEqual(ncur) && fparent.IsEqual(nparent) {
			return fcur, fparent
		}
		t.unlockNode(fparent, true)
		t.unlockNode(fcur, false)
	}
}

func (extenalTree *OptimisticLockTree[K, V]) Find(key K) (V, error) {
	// fmt.Printf("Find\n")
	t := extenalTree.mtree
	var nilValue V
	n, par := extenalTree.subFind(t.tree.root, nil, key)
	defer t.unlockNode(n, false)
	defer t.unlockNode(par, true)

	if n.IsNil() {
		return nilValue, btree.ErrorNodeNotFound
	} else {
		return n.GetValue(), nil
	}
}

func (extenalTree *OptimisticLockTree[K, V]) Insert(key K, value V) error {
	// fmt.Printf("Insert\n")
	t := extenalTree.mtree

	newNode := newMutexNode(key, value)
	t.lockNode(newNode, false)

	n, par, err := t.tree.NodeInsert(t.tree.root, nil, newNode, extenalTree.subFind)

	t.unlockNode(n, false)
	t.unlockNode(par, true)
	return err
}

func (extenalTree *OptimisticLockTree[K, V]) insertForDelete(lchild *mutexNode[K, V], parent *mutexNode[K, V], rchild *mutexNode[K, V]) {
	t := extenalTree.mtree

	t.lockNode(lchild, false)
	t.lockNode(rchild, false)
	key := parent.GetKey()
	var helper (func(*mutexNode[K, V], *mutexNode[K, V], *mutexNode[K, V]) (*mutexNode[K, V], *mutexNode[K, V]))

	canUnlockGrand := false
	helper = func(cur *mutexNode[K, V], parent *mutexNode[K, V], grandpar *mutexNode[K, V]) (*mutexNode[K, V], *mutexNode[K, V]) {
		t.lockNode(cur, false)

		if canUnlockGrand {
			t.unlockNode(grandpar, true)
		} else {
			if parent.IsEqual(lchild) {
				canUnlockGrand = true
			}
		}
		return t.tree.StepFind(cur, parent, key, helper)
	}

	nd, par := helper(lchild.GetRight(), lchild, parent)
	if !nd.IsNil() {
		log.Fatal("Get the same node during delete")
	} else {
		t.tree.NodeReplace(par, rchild)
	}

	t.unlockNode(par, true)
	t.unlockNode(rchild, false)
}

//	func (extenalTree *OptimisticLockTree[K, V]) subFindForDelete(cur *mutexNode[K, V], parent *mutexNode[K, V], key K) (*mutexNode[K, V], *mutexNode[K, V]) {
//		t := extenalTree.mtree
//		n, par := extenalTree.subFind(cur, parent, key)
//		t.lockNode(n.GetLeft(), false)
//		t.lockNode(n.GetRight(), false)
//		return n, par
//	}
func (extenalTree *OptimisticLockTree[K, V]) Delete(key K) error {
	// fmt.Printf("Delete\n")
	t := extenalTree.mtree

	_, par, err := t.tree.NodeDelete(t.tree.root, nil, key, extenalTree.subFind, extenalTree.insertForDelete)

	t.unlockNode(par, true)
	return err
}

func (extenalTree *OptimisticLockTree[K, V]) Print() {
	t := extenalTree.mtree
	cur := make(chan *mutexNode[K, V], 300)
	next := make(chan *mutexNode[K, V], 300)

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
		next = make(chan *mutexNode[K, V], 300)
		i++
		fmt.Printf("\n")
	}
}
