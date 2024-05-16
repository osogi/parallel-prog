package lincheck

import (
	"parallel-prog/4/btree"

	"golang.org/x/exp/constraints"
)

type SimpleTree[K constraints.Ordered, V any] struct {
	core map[K]V
}

func NewSimpleTree[K constraints.Ordered, V any]() *SimpleTree[K, V] {
	return &SimpleTree[K, V]{core: make(map[K]V)}
}

func (t *SimpleTree[K, V]) Insert(k K, v V) error {
	_, ok := t.core[k]

	if ok {
		return btree.ErrorSameKey
	} else {
		t.core[k] = v
		return nil
	}
}

func (t *SimpleTree[K, V]) Delete(k K) error {
	_, ok := t.core[k]

	if ok {
		delete(t.core, k)
		return nil
	} else {
		return btree.ErrorNodeNotFound
	}
}

func (t *SimpleTree[K, V]) Find(k K) (V, error) {
	var nilV V
	v, ok := t.core[k]

	if ok {
		return v, nil
	} else {
		return nilV, btree.ErrorNodeNotFound
	}
}
