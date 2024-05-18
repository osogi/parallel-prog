package trees_test

import "parallel-prog/4/btree"

func createNiceTree(tree btree.Btree[int, int], elems []int) {
	ln := len(elems)
	if ln != 0 {
		num := ln / 2
		tree.Insert(elems[num], elems[num])
		createNiceTree(tree, elems[0:num])
		createNiceTree(tree, elems[num+1:])
	}
}

func createRange(from int, to int) []int {
	sl := make([]int, to-from)
	c := 0
	for i := from; i < to; i++ {
		sl[c] = i
		c++
	}
	return sl
}
