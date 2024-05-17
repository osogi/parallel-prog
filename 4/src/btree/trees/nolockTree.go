package trees

import (
	"parallel-prog/4/btree"

	"golang.org/x/exp/constraints"
)

type NoLockTree[K constraints.Ordered, V any] struct {
	tree *CpsTree[K, V, *commonNode[K, V]]
}

func NewNoLockTree[K constraints.Ordered, V any]() *NoLockTree[K, V] {
	return &NoLockTree[K, V]{tree: NewCpsTree[K, V, *commonNode[K, V]](nil)}
}

func (t *NoLockTree[K, V]) subFind(cur *commonNode[K, V], parent *commonNode[K, V], key K) (*commonNode[K, V], *commonNode[K, V]) {
	return t.tree.NodeFind(cur, parent, nil, key, emptyFunc)
}
func (t *NoLockTree[K, V]) subInsert(cur *commonNode[K, V], parent *commonNode[K, V], newNode *commonNode[K, V]) {
	t.tree.NodeInsert(cur, parent, newNode, t.subFind)
}

func (t *NoLockTree[K, V]) Find(key K) (V, error) {
	var nilValue V

	n, _ := t.subFind(t.tree.root, nil, key)
	if n.IsNil() {
		return nilValue, btree.ErrorNodeNotFound
	} else {
		return n.GetValue(), nil
	}
}

func (t *NoLockTree[K, V]) Insert(key K, value V) error {
	newNode := newCommonNode(key, value)
	_, _, err := t.tree.NodeInsert(t.tree.root, nil, newNode, t.subFind)
	return err
}

func (t *NoLockTree[K, V]) Delete(key K) error {
	_, _, err := t.tree.NodeDelete(t.tree.root, nil, key, t.subFind, t.subInsert)
	return err
}
