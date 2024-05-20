package trees_test

import (
	"math/rand"
	"parallel-prog/4/btree"
	"parallel-prog/4/btree/trees"
	"sync"
	"testing"
)

var values []int
var numberSize = 2

func init() {
	rnd := rand.New(rand.NewSource(52))
	values = make([]int, numberSize)
	for i := range values {
		values[i] = rnd.Intn(1000)
	}
}

func TestFuzzing(t *testing.T) {
	trees := []struct {
		name string
		bst  btree.Btree[int, int]
	}{
		{name: "CoarseGrained", bst: trees.NewHardLockTree[int, int]()},
	}
	var wg sync.WaitGroup

	for _, tree := range trees {
		t.Run("Insert"+tree.name, func(t *testing.T) {
			for _, v := range values {
				wg.Add(1)
				go func(val int) {
					defer wg.Done()
					tree.bst.Insert(val, val)
				}(v)
			}
			wg.Wait()
			tree.bst.Print()

			// slog.Info(tree.bst.String())
			// if !tree.bst.Valid() {
			// 	t.Error("invalid tree")
			// }
		})

		t.Run("InsertDelete"+tree.name, func(t *testing.T) {
			for _, v := range values {
				wg.Add(1)
				go func(val int) {
					defer wg.Done()
					tree.bst.Insert(val, val)
				}(v)
			}
			wg.Wait()
			// slog.Info(tree.bst.String())

			for _, v := range values {
				wg.Add(1)
				go func(val int) {
					defer wg.Done()
					tree.bst.Delete(val)
				}(v)
			}
			wg.Wait()
			tree.bst.Print()
			// if !tree.bst.Valid() {
			// 	t.Error("invalid tree")
			// }
		})
	}
}
