package btree

import "golang.org/x/exp/constraints"

type Btree[K constraints.Ordered, V any] interface {
	Insert(K, V) error
	Delete(K) error
	Find(K) (V, error)
}
