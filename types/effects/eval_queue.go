package effects

// evalQueue represents a type-aligned sequence of Kleisli arrows.
type evalQueue[E any, A any, B any] interface {
	evalQ(eff E, input A) B // marker method - do not call
}

func (_ leafQ[E, A, B]) evalQ(eff E, input A) B    { panic("marker method") }
func (_ nodeQ[E, A, X, B]) evalQ(eff E, input A) B { panic("marker method") }

type leafQ[E any, A any, B any] struct {
	lifted func(arg A) Eff[E, B]
}

type nodeQ[E any, A any, X any, B any] struct {
	left  evalQueue[E, A, X]
	right evalQueue[E, X, B]
}

func liftQ[E any, A any, B any](f func(arg A) Eff[E, B]) evalQueue[E, A, B] {
	return leafQ[E, A, B]{lifted: f}
}

func concatQ[E any, A any, X any, B any](a2x evalQueue[E, A, X], x2b evalQueue[E, X, B]) evalQueue[E, A, B] {
	return nodeQ[E, A, X, B]{left: a2x, right: x2b}
}
