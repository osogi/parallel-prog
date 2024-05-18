package trees_test

import (
	"math/rand"
	"parallel-prog/4/btree"
	"parallel-prog/4/btree/lincheck"
	"parallel-prog/4/btree/trees"
	"testing"
	"time"
)

func RunLincheck(t *testing.T, emptyGen func() btree.Btree[int, int], threadsNum int, threadLen int, seed int64, diffSeedsNum int, repeatRunNum int) {
	rg := rand.New(rand.NewSource(seed))

	emptyCheck := func() btree.Btree[int, int] {
		return lincheck.NewSimpleTree[int, int]()
	}
	for i := 0; i < diffSeedsNum; i++ {
		// fmt.Println(i)
		c := lincheck.MakeChecker(emptyGen, emptyCheck, threadsNum, threadLen, int64(rg.Int()), 5*time.Second)
		err := c.RunCheck(repeatRunNum)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	}
}

func RunLincheckHard(t *testing.T, threadsNum int, threadLen int, seed int64, diffSeedsNum int, repeatRunNum int) {
	emptyGen := func() btree.Btree[int, int] {
		return trees.NewHardLockTree[int, int]()
	}
	RunLincheck(t, emptyGen, threadsNum, threadLen, seed, diffSeedsNum, repeatRunNum)
}

func RunLincheckFineGrade(t *testing.T, threadsNum int, threadLen int, seed int64, diffSeedsNum int, repeatRunNum int) {
	emptyGen := func() btree.Btree[int, int] {
		return trees.NewFineGradeLockTree[int, int]()
	}
	RunLincheck(t, emptyGen, threadsNum, threadLen, seed, diffSeedsNum, repeatRunNum)
}

func RunLincheckOptimistic(t *testing.T, threadsNum int, threadLen int, seed int64, diffSeedsNum int, repeatRunNum int) {
	emptyGen := func() btree.Btree[int, int] {
		return trees.NewOptimisticLockTree[int, int]()
	}
	RunLincheck(t, emptyGen, threadsNum, threadLen, seed, diffSeedsNum, repeatRunNum)
}

func TestLincheckHard_5threads_3ops(t *testing.T) {
	RunLincheckHard(t, 5, 3, 0xaaaa, 100, 10)
}

func TestLincheckFineGrade_5threads_3ops(t *testing.T) {
	RunLincheckFineGrade(t, 5, 3, 0xaaaa, 100, 10)
}

func TestLincheckFineGrade_3threads_7ops(t *testing.T) {
	RunLincheckFineGrade(t, 3, 7, 0xaaaa, 100, 10)
}

func TestLincheckFineGrade_2threads_20ops(t *testing.T) {
	RunLincheckFineGrade(t, 2, 20, 0xaaaa, 100, 10)
}

func TestLincheckOptimistic_5threads_3ops(t *testing.T) {
	RunLincheckOptimistic(t, 5, 3, 0xaaaa, 100, 10)
}

func TestLincheckOptimistic_3threads_7ops(t *testing.T) {
	RunLincheckOptimistic(t, 3, 7, 0xaaaa, 100, 10)
}

func TestLincheckOptimistic_2threads_20ops(t *testing.T) {
	RunLincheckOptimistic(t, 2, 20, 0xaaaa, 100, 10)
}
