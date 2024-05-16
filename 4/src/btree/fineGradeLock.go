package btree

// import (
// 	"fmt"
// 	"sync"

// 	"golang.org/x/exp/constraints"
// )

// type mutexNode[K constraints.Ordered, V any] struct {
// 	Key         K
// 	Value       V
// 	mutex       sync.Mutex
// 	Right, Left *mutexNode[K, V]
// }

// func newMutexNode[K constraints.Ordered, V any](key K, value V) *mutexNode[K, V] {
// 	return &mutexNode[K, V]{Key: key, Value: value, mutex: sync.Mutex{}, Right: nil, Left: nil}
// }

// func (p *mutexNode[K, V]) insert(n *mutexNode[K, V]) error {
// 	errorNotNil := fmt.Errorf("can't insert node instead not nil node")
// 	if n.Key == p.Key {
// 		return fmt.Errorf("can't insert node to parent with same key")
// 	} else if n.Key < p.Key {
// 		if p.Left != nil {
// 			return errorNotNil
// 		} else {
// 			p.Left = n
// 		}
// 	} else {
// 		if p.Right != nil {
// 			return errorNotNil
// 		} else {
// 			p.Right = n
// 		}
// 	}

// 	return nil
// }

// type SoftLockTree[K constraints.Ordered, V any] struct {
// 	root *mutexNode[K, V]
// }

// func NewSoftLockTree[K constraints.Ordered, V any]() *HardLockTree[K, V] {
// 	return &HardLockTree[K, V]{root: nil, mutex: sync.Mutex{}}
// }
