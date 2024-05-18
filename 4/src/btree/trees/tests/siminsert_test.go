package trees_test

import (
	"parallel-prog/4/btree"
	"parallel-prog/4/btree/trees"
	"sync"
	"sync/atomic"
	"testing"
)

/*
This test checks that tree.Insert() works correctly when called
from multiple threads on the same element, i.e. that the element is added exactly once.

Considering that the testSimultaneousInsert function calls subTestSimultaneousInsert several times in a row,
it also checks that the tree remains "correct" for future operations
*/
func subTestSimultaneousInsert(t *testing.T, tree btree.Btree[int, int], insertNum int, threadNum int) {
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
			err := tree.Insert(insertNum, insertNum)
			if err != nil {
				if err != btree.ErrorSameKey {
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
		t.Errorf("Insert one element several times (%d)", c.Load())
	}
}

func testSimultaneousInsert(t *testing.T, tree btree.Btree[int, int], insLst []int, threadNum int) {
	for _, ins := range insLst {
		subTestSimultaneousInsert(t, tree, ins, threadNum)
	}
}

func runTestSimultaneousInsert(t *testing.T, tree btree.Btree[int, int], treeLst []int, insLst []int, threadNum int) {
	createNiceTree(tree, treeLst)
	testSimultaneousInsert(t, tree, insLst, threadNum)
}

func runSimIns1(t *testing.T, emptyTree btree.Btree[int, int]) {
	runTestSimultaneousInsert(t, emptyTree, []int{}, []int{50, 12, 32, 16, 37, 6, 52}, 100)
}

func runSimIns2(t *testing.T, emptyTree btree.Btree[int, int]) {
	sl := append(createRange(0, 25), createRange(75, 100)...)
	runTestSimultaneousInsert(t, emptyTree, sl, createRange(25, 75), 100)
}

// simultaneous root insert
func runSimIns3(t *testing.T, emptyTree btree.Btree[int, int]) {
	for i := 0; i < 100; i++ {
		runTestSimultaneousInsert(t, emptyTree, []int{}, []int{2}, 100)
		emptyTree.Delete(2)
	}
}

func TestSimInsertHardLock(t *testing.T) {
	runSimIns1(t, trees.NewHardLockTree[int, int]())
	runSimIns2(t, trees.NewHardLockTree[int, int]())
	runSimIns3(t, trees.NewHardLockTree[int, int]())
}

func TestSimInsertFineGrade(t *testing.T) {
	runSimIns1(t, trees.NewFineGradeLockTree[int, int]())
	runSimIns2(t, trees.NewFineGradeLockTree[int, int]())
	runSimIns3(t, trees.NewFineGradeLockTree[int, int]())
}

func TestSimInsertOptimistic(t *testing.T) {
	runSimIns1(t, trees.NewOptimisticLockTree[int, int]())
	runSimIns2(t, trees.NewOptimisticLockTree[int, int]())
	runSimIns3(t, trees.NewOptimisticLockTree[int, int]())
}
