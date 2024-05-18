package trees_test

import (
	"parallel-prog/4/btree"
	"parallel-prog/4/btree/trees"
	"sync"
	"sync/atomic"
	"testing"
)

/*
This test checks that tree.Delete() works correctly when called
from multiple threads on the same element, i.e. that the element is removed exactly once.

Considering that the testSimultaneousDelete function calls subTestSimultaneousDelete several times in a row,
it also checks that the tree remains "correct" for future operations
*/
func subTestSimultaneousDelete(t *testing.T, tree btree.Btree[int, int], del int, threadNum int) {
	randevu := sync.WaitGroup{}
	randevu.Add(threadNum)
	end := sync.WaitGroup{}
	end.Add(threadNum)

	var c atomic.Int32
	c.Store(0)
	for i := 0; i < threadNum; i++ {
		go func() {
			randevu.Done()
			randevu.Wait()
			err := tree.Delete(del)
			if err != nil {
				if err != btree.ErrorNodeNotFound {
					t.Errorf("Can't delete existing element, because:\n\t%s", err.Error())
				}
			} else {
				c.Add(1)
			}
			end.Done()
		}()
	}
	end.Wait()
	if c.Load() != 1 {
		t.Errorf("Delete one element several times (%d)", c.Load())
	}
}

func testSimultaneousDelete(t *testing.T, tree btree.Btree[int, int], delLst []int, threadNum int) {
	for _, del := range delLst {
		subTestSimultaneousDelete(t, tree, del, threadNum)
	}
}

func runTestSimultaneousDelete(t *testing.T, tree btree.Btree[int, int], treeLst []int, delLst []int, threadNum int) {
	createNiceTree(tree, treeLst)
	testSimultaneousDelete(t, tree, delLst, threadNum)
}

func runSimDel1(t *testing.T, emptyTree btree.Btree[int, int]) {
	runTestSimultaneousDelete(t, emptyTree, createRange(0, 100), []int{50, 12, 32, 16, 37, 6, 52}, 100)

}

func runSimDel2(t *testing.T, emptyTree btree.Btree[int, int]) {
	runTestSimultaneousDelete(t, emptyTree, createRange(0, 100), createRange(0, 100), 100)
}

// simultaneous root delete
func runSimDel3(t *testing.T, emptyTree btree.Btree[int, int]) {
	for i := 0; i < 100; i++ {
		emptyTree.Insert(2, 2)
		runTestSimultaneousDelete(t, emptyTree, []int{}, []int{2}, 100)
	}

}

func TestSimDeleteHardLock(t *testing.T) {
	runSimDel1(t, trees.NewHardLockTree[int, int]())
	runSimDel2(t, trees.NewHardLockTree[int, int]())
	runSimDel3(t, trees.NewHardLockTree[int, int]())
}

func TestSimDeleteFineGrade(t *testing.T) {
	runSimDel1(t, trees.NewFineGradeLockTree[int, int]())
	runSimDel2(t, trees.NewFineGradeLockTree[int, int]())
	runSimDel3(t, trees.NewFineGradeLockTree[int, int]())
}

func TestSimDeleteOptimistic(t *testing.T) {
	runSimDel1(t, trees.NewOptimisticLockTree[int, int]())
	runSimDel2(t, trees.NewOptimisticLockTree[int, int]())
	runSimDel3(t, trees.NewOptimisticLockTree[int, int]())
}
