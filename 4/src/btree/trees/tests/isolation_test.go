package trees_test

import (
	"parallel-prog/4/btree"
	"parallel-prog/4/btree/trees"
	"sync"
	"testing"
)

/*
This test checks the "isolation" of elements,
i.e. that when working (inserting, deleting) node x,
other nodes remain accessible and do not disappear from the tree.

keys - must contain the keys with which the test will interact. If value keys[k] equal:
  - true: this key will first be deleted and then inserted
  - false: this key will be checked to still in the tree using the tree.Find() method
*/
func testIsolation(t *testing.T, tree btree.Btree[int, int], keys map[int]bool, operationCount int) {
	randevu := sync.WaitGroup{}
	randevu.Add(len(keys))
	end := sync.WaitGroup{}
	end.Add(len(keys))

	for k, del := range keys {
		del, k := del, k
		go func() {
			randevu.Done()
			randevu.Wait()
			for i := 0; i < operationCount; i++ {
				if del {
					err := tree.Delete(k)
					if err != nil {
						t.Errorf("Can't delete existing element, because:\n\t%s", err.Error())
					}
					err = tree.Insert(k, k)
					if err != nil {
						t.Errorf("Can't insert existing element, because:\n\t%s", err.Error())
					}
				} else {
					_, err := tree.Find(k)
					if err != nil {
						t.Errorf("Can't find existing element, because:\n\t%s", err.Error())
					}
				}
			}
			end.Done()
		}()
	}
	end.Wait()
}

func runTestIsolation(t *testing.T, tree btree.Btree[int, int], treeLst []int, delLst []int, operationCount int) {
	createNiceTree(tree, treeLst)
	del := make(map[int]bool)

	for _, c := range treeLst {
		del[c] = false
	}
	for _, c := range delLst {
		del[c] = true
	}
	testIsolation(t, tree, del, operationCount)
}

func runIsol1(t *testing.T, emptyTree btree.Btree[int, int]) {
	runTestIsolation(t, emptyTree, createRange(0, 100), []int{50, 12, 32, 16, 37, 6, 52}, 1000)
}

func runIsol2(t *testing.T, emptyTree btree.Btree[int, int]) {
	runTestIsolation(t, emptyTree, createRange(0, 100), createRange(25, 75), 1000)
}

func runIsol3(t *testing.T, emptyTree btree.Btree[int, int]) {
	runTestIsolation(t, emptyTree, createRange(0, 16), createRange(0, 15), 1000)
}

func TestIsolationHardLock(t *testing.T) {
	runIsol1(t, trees.NewHardLockTree[int, int]())
	runIsol2(t, trees.NewHardLockTree[int, int]())
	runIsol3(t, trees.NewHardLockTree[int, int]())
}

func TestIsolationFineGrade(t *testing.T) {
	runIsol1(t, trees.NewFineGradeLockTree[int, int]())
	runIsol2(t, trees.NewFineGradeLockTree[int, int]())
	runIsol3(t, trees.NewFineGradeLockTree[int, int]())
}

func TestIsolationOptimistic(t *testing.T) {
	runIsol1(t, trees.NewOptimisticLockTree[int, int]())
	runIsol2(t, trees.NewOptimisticLockTree[int, int]())
	runIsol3(t, trees.NewOptimisticLockTree[int, int]())
}
