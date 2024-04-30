package btree

import "errors"

var ErrorSameKey = errors.New("node with same key already exists in the tree")
var ErrorNodeNotFound = errors.New("node with this key not found")
