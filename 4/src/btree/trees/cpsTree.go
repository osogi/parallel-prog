package trees

import (
	"parallel-prog/4/btree"

	"golang.org/x/exp/constraints"
)

type CpsTree[K constraints.Ordered, V any, NT Node[K, V, NT]] struct {
	root NT
}

func NewCpsTree[K constraints.Ordered, V any, NT Node[K, V, NT]](root NT) *CpsTree[K, V, NT] {
	return &CpsTree[K, V, NT]{root: root}
}

func (t *CpsTree[K, V, NT]) StepFind(cur NT, parent NT, key K, cont func(NT, NT, NT) (NT, NT)) (NT, NT) {
	if !cur.IsNil() {
		curkey := cur.GetKey()
		if curkey == key {
			return cur, parent
		} else if key < curkey {
			return cont(cur.GetLeft(), cur, parent)
		} else {
			return cont(cur.GetRight(), cur, parent)
		}
	}
	return cur, parent
}

// func (t *CpsTree[K, V, NT]) StepError(cur NT, parent NT, key K, cont func(NT, NT, NT) error) error {
// 	if !cur.IsNil() {
// 		curkey := cur.GetKey()
// 		if curkey == key {
// 			return ErrorSameKey
// 		} else if key < curkey {
// 			return cont(cur.GetLeft(), cur, parent)
// 		} else {
// 			return cont(cur.GetRight(), cur, parent)
// 		}
// 	}
// 	return ErrorNilNode
// }

type FindFunType[K any, NT any] func(NT, NT, K) (NT, NT)

func (t *CpsTree[K, V, NT]) NodeFind(cur NT, parent NT, grandpar NT, key K, preStep func(NT, NT, NT)) (NT, NT) {
	var helper (func(NT, NT, NT) (NT, NT))

	helper = func(child NT, cur NT, parent NT) (NT, NT) {
		preStep(child, cur, parent)
		return t.StepFind(child, cur, key, helper)
	}

	return helper(cur, parent, grandpar)
}

func (t *CpsTree[K, V, NT]) NodeReplace(parent NT, newNode NT) NT {
	if parent.IsNil() {
		t.root = newNode
	} else {
		if newNode.GetKey() < parent.GetKey() {
			parent.setLeft(newNode)
		} else {
			parent.setRight(newNode)
		}
	}
	return newNode
}

func (t *CpsTree[K, V, NT]) NodeReplacePrev(parent NT, prev NT, newNode NT) NT {
	if parent.IsNil() {
		t.root = newNode
	} else {
		if prev.IsEqual(parent.GetLeft()) {
			parent.setLeft(newNode)
		} else {
			parent.setRight(newNode)
		}
	}
	return newNode
}

func (t *CpsTree[K, V, NT]) NodeInsert(cur NT, parent NT, newNode NT, findFun FindFunType[K, NT]) (NT, NT, error) {
	key := newNode.GetKey()
	finded, fparent := findFun(cur, parent, key)
	var err error = nil
	if !finded.IsNil() {
		err = btree.ErrorSameKey
	} else {
		finded = t.NodeReplace(fparent, newNode)
	}
	return finded, fparent, err
}

// func (t *CpsTree[K, V, NT]) NodeInsert(cur NT, parent NT, newNode NT, preStep func(NT, NT, NT), preInsert func(NT, NT, NT), postInsert func(NT, NT, NT)) error {
// 	var helper (func(NT, NT, NT) error)
// 	var nilNode NT
// 	key := newNode.GetKey()

// 	helper = func(child NT, cur NT, parent NT) error {
// 		preStep(child, cur, parent)
// 		if child.IsNil() {
// 			preInsert(newNode, cur, parent)
// 			if cur.IsNil() {
// 				t.root = newNode
// 			} else if key < cur.GetKey() {
// 				cur.setLeft(newNode)
// 			} else {
// 				cur.setRight(newNode)
// 			}
// 			postInsert(newNode, cur, parent)
// 			return nil
// 		} else {
// 			return t.StepError(child, cur, key, helper)
// 		}
// 	}

// 	return helper(cur, parent, nilNode)
// }

func (t *CpsTree[K, V, NT]) ZeroOneChildNodeDelete(target NT, parent NT) NT {
	tl := target.GetLeft()
	tr := target.GetRight()

	if tl.IsNil() {
		return t.NodeReplacePrev(parent, target, tr)
	} else if tr.IsNil() {
		return t.NodeReplacePrev(parent, target, tl)
	}
	return target
}

func (t *CpsTree[K, V, NT]) NodeDelete(cur NT, parent NT, key K, findFun FindFunType[K, NT], insertFun func(NT, NT, NT)) (NT, NT, error) {
	finded, fparent := findFun(cur, parent, key)
	var err error = nil
	var newChild NT = finded

	if finded.IsNil() {
		err = btree.ErrorNodeNotFound
	} else {
		newChild = t.ZeroOneChildNodeDelete(finded, fparent)
		if !newChild.IsEqual(finded) {
			err = nil
		} else {
			fl := finded.GetLeft()
			fr := finded.GetRight()
			insertFun(fl, finded, fr)
			newChild = t.NodeReplacePrev(fparent, finded, fl)
		}
	}
	return newChild, fparent, err
}

// func (t *CpsTree[K, V, NT]) DeleteTargetNode(target NT, parent NT, preStep func(NT, NT, NT), preDelete func(NT, NT, NT), postDelete func(NT, NT, NT)) error {
// 	if target.IsNil() {
// 		return ErrorNodeNotFound
// 	} else {
// 		tl := target.GetLeft()
// 		tr := target.GetRight()

// 		if tl.IsNil() {
// 			preDelete(parent, target, tr)
// 			t.NodeReplace(parent, target, tr)
// 			postDelete(parent, target, tr)
// 		} else if tr.IsNil() {
// 			preDelete(parent, target, tl)
// 			t.NodeReplace(parent, target, tl)
// 			postDelete(parent, target, tl)
// 		} else {
// 			npar := target
// 			ncur := target.GetLeft()
// 			nchild := ncur.GetRight()

// 			preStep(ncur, npar, nchild)
// 			for !nchild.IsNil() {
// 				npar = ncur
// 				ncur = nchild
// 				nchild = nchild.GetRight()
// 				preStep(ncur, npar, nchild)
// 			}

// 			err := t.DeleteTargetNode(ncur, npar, preDelete, preStep, postDelete)
// 			if err != nil {
// 				return err
// 			}

// 			preDelete(parent, target, ncur)
// 			ncur.setRight(target.GetRight())
// 			ncur.setLeft(target.GetLeft())
// 			t.NodeReplace(parent, target, ncur)
// 			postDelete(parent, target, ncur)

// 		}
// 	}

// 	return nil
// }

// func (t *CpsTree[K, V, NT]) NodeDelete(cur NT, parent NT, key K, preStep func(NT, NT, NT), additionalPreStep func(NT, NT, NT), preDelete func(NT, NT, NT), postDelete func(NT, NT, NT)) error {
// 	var helper (func(NT, NT, NT) error)
// 	var nilNode NT

// 	helper = func(child NT, cur NT, parent NT) error {
// 		preStep(child, cur, parent)
// 		if child.IsNil() {
// 			return ErrorNodeNotFound
// 		}
// 		if child.GetKey() == key {
// 			return t.DeleteTargetNode(child, cur, additionalPreStep, preDelete, postDelete)
// 		} else {
// 			return t.StepError(child, cur, key, helper)
// 		}
// 	}

// 	return helper(cur, parent, nilNode)
// }
