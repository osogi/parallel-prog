package btree

import "golang.org/x/exp/constraints"

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

func (t *CpsTree[K, V, NT]) StepError(cur NT, parent NT, key K, cont func(NT, NT, NT) error) error {
	var nilNode NT

	if !cur.IsNil() {
		curkey := cur.GetKey()
		if curkey == key {
			return cont(nilNode, cur, parent)
		} else if key < curkey {
			return cont(cur.GetLeft(), cur, parent)
		} else {
			return cont(cur.GetRight(), cur, parent)
		}
	}
	return ErrorNilNode
}

func (t *CpsTree[K, V, NT]) NodeFind(cur NT, parent NT, key K, preStep func(NT, NT, NT)) (NT, NT) {
	var nilNode NT

	var helper (func(NT, NT, NT) (NT, NT))

	helper = func(child NT, cur NT, parent NT) (NT, NT) {
		preStep(child, cur, parent)
		return t.StepFind(cur, parent, key, helper)
	}

	return helper(cur, parent, nilNode)
}

func (t *CpsTree[K, V, NT]) NodeInsert(cur NT, parent NT, newNode NT, preStep func(NT, NT, NT), preInsert func(NT, NT, NT), postInsert func(NT, NT, NT)) error {
	var helper (func(NT, NT, NT) error)
	var nilNode NT
	key := newNode.GetKey()

	helper = func(child NT, cur NT, parent NT) error {
		preStep(child, cur, parent)
		if child.IsNil() {
			preInsert(newNode, cur, parent)
			if cur.IsNil() {
				t.root = newNode
			}
			if key < cur.GetKey() {
				cur.setLeft(newNode)
			} else {
				cur.setRight(newNode)
			}
			postInsert(newNode, cur, parent)
			return nil
		} else {
			return t.StepError(cur, parent, key, helper)
		}
	}

	return helper(cur, parent, nilNode)
}

func (t *CpsTree[K, V, NT]) NodeReplace(parent NT, prev NT, new NT) {
	if parent.IsNil() {
		t.root = new
	} else {
		if parent.GetLeft().IsEqual(prev) {
			parent.setLeft(new)
		} else {
			parent.setRight(new)
		}
	}
}

func (t *CpsTree[K, V, NT]) DeleteTargetNode(target NT, parent NT, preStep func(NT, NT, NT), preDelete func(NT, NT, NT), postDelete func(NT, NT, NT)) error {
	if target.IsNil() {
		return ErrorNodeNotFound
	} else {
		tl := target.GetLeft()
		tr := target.GetRight()

		if tl.IsNil() {
			preDelete(parent, target, tr)
			t.NodeReplace(parent, target, tr)
			postDelete(parent, target, tr)
		} else if tr.IsNil() {
			preDelete(parent, target, tl)
			t.NodeReplace(parent, target, tl)
			postDelete(parent, target, tl)
		} else {
			npar := target
			ncur := target.GetLeft()
			nchild := ncur.GetRight()

			preStep(ncur, npar, nchild)
			for !nchild.IsNil() {
				npar = ncur
				ncur = nchild
				nchild = nchild.GetRight()
				preStep(ncur, npar, nchild)
			}

			err := t.DeleteTargetNode(ncur, npar, preDelete, preStep, postDelete)
			if err != nil {
				return err
			}

			preDelete(parent, target, ncur)
			ncur.setRight(target.GetRight())
			t.NodeReplace(parent, target, ncur)
			postDelete(parent, target, ncur)

		}
	}

	return nil
}

func (t *CpsTree[K, V, NT]) NodeDelete(cur NT, parent NT, key K, preStep func(NT, NT, NT), additionalPreStep func(NT, NT, NT), preDelete func(NT, NT, NT), postDelete func(NT, NT, NT)) error {
	var helper (func(NT, NT, NT) error)
	var nilNode NT

	helper = func(child NT, cur NT, parent NT) error {
		preStep(child, cur, parent)
		if child.IsNil() {
			return ErrorNodeNotFound
		}
		if child.GetKey() == key {
			return t.DeleteTargetNode(child, cur, additionalPreStep, preDelete, postDelete)
		} else {
			return t.StepError(cur, parent, key, helper)
		}
	}

	return helper(cur, parent, nilNode)
}
