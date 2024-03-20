package main

import (
	"fmt"
	"parallel-prog/1/stack"
)

func main() {
	fmt.Println("So you have a mother too!")
	var a *stack.CommonStack[int] // stack.NewCommonStack[int]()
	err := a.Push(12)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
