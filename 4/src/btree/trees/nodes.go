package trees

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type Node[K constraints.Ordered, V any, NT any] interface {
	GetKey() K
	setKey(K)

	GetValue() V
	setValue(V)

	GetRight() NT
	setRight(NT)

	GetLeft() NT
	setLeft(NT)

	IsNil() bool
	IsEqual(NT) bool
}

type commonNode[K constraints.Ordered, V any] struct {
	Key         K
	Value       V
	Right, Left *commonNode[K, V]
}

func newCommonNode[K constraints.Ordered, V any](key K, value V) *commonNode[K, V] {
	return &commonNode[K, V]{Key: key, Value: value, Right: nil, Left: nil}
}

func (n *commonNode[K, V]) GetKey() K {
	return n.Key
}
func (n *commonNode[K, V]) setKey(key K) {
	n.Key = key
}
func (n *commonNode[K, V]) GetValue() V {
	return n.Value
}
func (n *commonNode[K, V]) setValue(value V) {
	n.Value = value
}
func (n *commonNode[K, V]) GetRight() *commonNode[K, V] {
	return n.Right
}
func (n *commonNode[K, V]) setRight(right *commonNode[K, V]) {
	n.Right = right
}
func (n *commonNode[K, V]) GetLeft() *commonNode[K, V] {
	return n.Left
}
func (n *commonNode[K, V]) setLeft(left *commonNode[K, V]) {
	n.Left = left
}
func (n *commonNode[K, V]) IsNil() bool {
	return n == nil
}
func (n *commonNode[K, V]) IsEqual(other *commonNode[K, V]) bool {
	return n == other
}

type mutexNode[K constraints.Ordered, V any] struct {
	Key         K
	Value       V
	mutex       sync.Mutex
	Right, Left *mutexNode[K, V]
}

func newMutexNode[K constraints.Ordered, V any](key K, value V) *mutexNode[K, V] {
	return &mutexNode[K, V]{Key: key, Value: value, mutex: sync.Mutex{}, Right: nil, Left: nil}
}
func (n *mutexNode[K, V]) GetKey() K {
	return n.Key
}
func (n *mutexNode[K, V]) setKey(key K) {
	n.Key = key
}
func (n *mutexNode[K, V]) GetValue() V {
	return n.Value
}
func (n *mutexNode[K, V]) setValue(value V) {
	n.Value = value
}
func (n *mutexNode[K, V]) GetRight() *mutexNode[K, V] {
	return n.Right
}
func (n *mutexNode[K, V]) setRight(right *mutexNode[K, V]) {
	n.Right = right
}
func (n *mutexNode[K, V]) GetLeft() *mutexNode[K, V] {
	return n.Left
}
func (n *mutexNode[K, V]) setLeft(left *mutexNode[K, V]) {
	n.Left = left
}
func (n *mutexNode[K, V]) IsNil() bool {
	return n == nil
}
func (n *mutexNode[K, V]) IsEqual(other *mutexNode[K, V]) bool {
	return n == other
}
