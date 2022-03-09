package functional

type Foldable[A any] interface {
	FoldLeft(fn FoldLeftFn[A])
}

type FoldLeftFn[A any] interface {
	Next(elem A)
}

func FoldLeft[FA Foldable[A], A any, B any](foldable FA, initial B, foldFn func(elem A, state B) B) B {
	fi := &foldImpl[A, B]{initial, foldFn}
	foldable.FoldLeft(fi)
	return fi.Value()
}

type foldImpl[A any, S any] struct {
	state S
	next  func(elem A, state S) S
}

func (f *foldImpl[_, S]) Value() S {
	return f.state
}

func (f *foldImpl[A, S]) Next(elem A) {
	f.state = f.next(elem, f.state)
}
