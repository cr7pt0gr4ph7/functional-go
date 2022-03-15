package effects 

type Arr[E any, A any, B any] func(arg A) Eff[E, B]

// evalQueue represents a type-aligned sequence of Kleisli arrows.
type evalQueue[E any, A any, B any] interface {
	evalQ(eff E, input A) B // marker method - do not call
	evalQBase
	evalQIn[E, A]
	evalQOut[E, B]
}

type evalQBase interface {
	isLeaf() bool
	apply(start any) effBase
	leftBase() evalQBase
	rightBase() evalQBase
	qApply(start any) effBase
	qPrepend(effect EffectTag, queue evalQBase) effBase
}

type evalQIn[E any, A any] interface {
	evalQBase
}

type evalQOut[E any, B any] interface {
	evalQBase
}

type evalQ[E any, A any, B any] interface {
	evalQIn[E, A]
	evalQOut[E, B]
}

func (_ leafQ[E, A, B]) evalQ(eff E, input A) B    { panic("marker method") }
func (_ nodeQ[E, A, X, B]) evalQ(eff E, input A) B { panic("marker method") }

func (l leafQ[E, A, B]) isLeaf() bool            { return true }
func (l leafQ[E, A, B]) apply(start any) effBase { return l.lifted(start.(A)) }
func (l leafQ[E, A, B]) leftBase() evalQBase     { return nil }
func (l leafQ[E, A, B]) rightBase() evalQBase    { return nil }

func (l leafQ[E, A, B]) qApply(start any) effBase {
	return qApply(start.(A), evalQueue[E, A, B](l))
}

func (l leafQ[E, A, B]) qPrepend(effect EffectTag, queue evalQBase) effBase {
	return newCont[E, Start, B](Union[E, Start](effect), concatQ[E, Start, A, B](queue.(evalQueue[E, Start, A]), l))
}

func (t nodeQ[E, A, X, B]) isLeaf() bool            { return false }
func (t nodeQ[E, A, X, B]) apply(start any) effBase { panic("not a leaf node") }
func (t nodeQ[E, A, X, B]) leftBase() evalQBase     { return t.left }
func (t nodeQ[E, A, X, B]) rightBase() evalQBase    { return t.right }

func (t nodeQ[E, A, X, B]) qApply(start any) effBase {
	return qApply(start.(A), evalQueue[E, A, B](t))
}

func (t nodeQ[E, A, X, B]) qPrepend(effect EffectTag, queue evalQBase) effBase {
	return newCont[E, Start, B](Union[E, Start](effect), concatQ[E, Start, A, B](queue.(evalQueue[E, Start, A]), t))
}

func (t nodeQ2[E, B]) isLeaf() bool            { return false }
func (t nodeQ2[E, B]) apply(start any) effBase { panic("not a leaf node") }
func (t nodeQ2[E, B]) leftBase() evalQBase     { return t.left }
func (t nodeQ2[E, B]) rightBase() evalQBase    { return t.right }

func (t nodeQ2[E, B]) qApply(start any) effBase {
	return qApply2[E, B](start, evalQOut[E, B](t))
}

func (t nodeQ2[E, B]) qPrepend(effect EffectTag, queue evalQBase) effBase {
	panic("types")
	// return newCont(m.effect, concatQ(m.queue, t))
}

type leafQ[E any, A any, B any] struct {
	lifted Arr[E, A, B]
}

type nodeQ[E any, A any, X any, B any] struct {
	left  evalQueue[E, A, X]
	right evalQueue[E, X, B]
}

type nodeQ2[E any, B any] struct {
	left  evalQBase
	right evalQBase
}

func liftQ[E any, A any, B any](f Arr[E, A, B]) evalQueue[E, A, B] {
	return leafQ[E, A, B]{lifted: f}
}

func concatQ[E any, A any, X any, B any](a2x evalQueue[E, A, X], x2b evalQueue[E, X, B]) evalQueue[E, A, B] {
	return nodeQ[E, A, X, B]{left: a2x, right: x2b}
}

func concatQ2Out[E any, B any](a2x evalQBase, x2b evalQOut[E, B]) evalQOut[E, B] {
	return nodeQ2[E, B]{left: a2x, right: x2b}
}

func qApply[E any, A any, B any](start A, q evalQueue[E, A, B]) Eff[E, B] {
	switch q := q.(type) {
	case leafQ[E, A, B]:
		return q.lifted(start)
	default:
		// q hasType nodeQ[E, A, X, B] where exists(X):
		return qApplyInner[E, A, B](
			start,         // A
			q.leftBase(),  // evalQueue[E, A, X]
			q.rightBase(), // evalQueue[E, X, B]
		) // => Eff[E, B]
	}
}

func qApply2[E any, B any](start any, q evalQOut[E, B]) Eff[E, B] {
	if q.isLeaf() {
		// start hasType A
		// q hasType leafQ[E, A, B]
		return q.apply(start).(Eff[E, B])
	} else {
		// q hasType nodeQ[E, A, X, B] where exists(X):
		return qApplyInner[E, any /* A */, B](
			start,         // A
			q.leftBase(),  // evalQueue[E, A, X]
			q.rightBase(), // evalQueue[E, X, B]
		) // => Eff[E, B]
	}
}

func qApplyInner[E any, A any, B any](start A, tl evalQIn[E, A], tr evalQOut[E, B]) Eff[E, B] {
	// (tl hasType evalQueue[E, A, X]
	//  tr hasType evalQueue[E, X, B]) where exists(X)

	if tl.isLeaf() {
		// tl hasType leafQ[E, A, X]
		// therefore: tl.lifted(start) hasType Eff[E, X]
		return qBind2[E, B]( // Eff[E, X] => evalQueue[E, X, B] => Eff[E, B]
			tl.apply(start), // (start: A) => (tl.apply: A => Eff[E, X]) => Eff[E, X]
			tr,              // evalQueue[E, X, B]
		) // => Eff[E, B]
	} else {
		// tl hasType nodeQ[E, A, Y, X] where exists(Y)
		// therefore tl.left  hasType evalQueue[E, A, Y]
		//           tl.right hasType evalQueue[E, Y, X]
		return qApplyInner[E, A, B](
			start,         // A
			tl.leftBase(), // evalQueue[E, A, Y]
			concatQ2Out[E, B](
				tl.rightBase(), // evalQueue[E, Y, X]
				tr,             // evalQueue[E, X, B]
			), // => evalQueue[E, Y, B]
		) // => Eff[E, B]
	}
}

func qCompose[E any, A any, B any, C any](a2b evalQueue[E, A, B], b2c func(eff Eff[E, B]) Eff[E, C]) Arr[E, A, C] {
	return func(start A) Eff[E, C] {
		return b2c( // Eff[E, B] => Eff[E, C]
			qApply(
				start, // A
				a2b,   // evalQueue[E, A, B]
			), // => Eff[E, B]
		) // => Eff[E, C]
	} // => (A => Eff[E, C])
}

func qBind[E any, A any, B any](e Eff[E, A], k evalQueue[E, A, B]) Eff[E, B] {
	switch m := e.EffImpl.(type) {
	case Pure[E, A]:
		return qApply(m.value, k)
	case Cont[E, Start, A]:
		return newCont(m.effect, concatQ(m.queue, k))
	default:
		panic("unreachable")
	}
}

func qBind2[E any, B any](e effBase, k evalQBase) Eff[E, B] {
	// (e hasType Eff[E, X]
	//  k hasType evalQueue[E, X, B]) where exists(X)

	if m := e.impl(); m.isPure() {
		// m hasType Pure[E, X]
		// k hasType evalQueue[E, X, B]
		// k.qApply(...) hasType Eff[E, B]
		return k.qApply(m.pureValue()).(Eff[E, B])
	} else {
		// m hasType Cont[E, Start, X]:
		// k hasType evalQueue[E, X, B]
		// k.qPrepend(...) hasType Eff[E, B]
		return k.qPrepend(m.getEffect(), m.getQueue()).(Eff[E, B])
	}
}
