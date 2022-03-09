package functional

import (
	"github.com/cr7pt0gr4ph7/functional-go/eval"
)

type Slice[A any] []A

func (s Slice[A]) FoldLeft(fn FoldLeftFn[A]) {
	for _, elem := range s {
		fn.Next(elem)
	}
}

func sliceFoldLeft[A any, B any](s Slice[A], initial B, foldFn func(elem A, state B) B) B {
	acc := initial
	for _, elem := range s {
		acc = foldFn(s, acc)
	}
	return acc
}

func sliceFoldRight[A any, B any](s Slice[A], final eval.Eval[B], foldFn func(elem A, state eval.Eval[B]) eval.Eval[B]) eval.Eval[B] {
	return eval.Defer(func() eval.Eval[B] {
		if len(s) == 0 {
			return final
		}
		return foldFn(s[0], sliceFoldRight(s[1:], final, foldFn))
	})
}
