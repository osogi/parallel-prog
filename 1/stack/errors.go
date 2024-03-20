package stack

import "errors"

var ErrorNilPointer = errors.New("nil pointer founded, expected pointer to CommonStack")
var ErrorEmptyStack = errors.New("stack is empty")
