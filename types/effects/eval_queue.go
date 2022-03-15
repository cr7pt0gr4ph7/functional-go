package effects 

type Arr[E any, A any, B any] func(arg A) Eff[E, B]

// evalQueue represents a type-aligned sequence of Kleisli arrows.
type evalQueue[E any, A any, B any] interface {
	evalQ(eff E, input A) B // marker method - do not call
	evalTreeNode
	evalLeftNode[E, A]
	evalRightNode[E, B]
}

// A tree node where neither the input type nor the output type
// are statically known at compile-time.
type evalTreeNode interface {
	isLeaf() bool
	leftTree() evalTreeNode
	rightTree() evalTreeNode
}

func (l leafQ[E, A, B]) isLeaf() bool                    { return true }
func (l leafQ[E, A, B]) leftTree() evalTreeNode          { return nil }
func (l leafQ[E, A, B]) leftTree_() evalLeftNode[E, B]   { return nil }
func (l leafQ[E, A, B]) rightTree() evalTreeNode         { return nil }
func (l leafQ[E, A, B]) rightTree_() evalRightNode[E, B] { return nil }

func (t nodeQ[E, A, X, B]) isLeaf() bool                    { return false }
func (t nodeQ[E, A, X, B]) leftTree() evalTreeNode          { return t.left }
func (t nodeQ[E, A, X, B]) leftTree_() evalLeftNode[E, A]   { return t.left }
func (t nodeQ[E, A, X, B]) rightTree() evalTreeNode         { return t.right }
func (t nodeQ[E, A, X, B]) rightTree_() evalRightNode[E, B] { return t.right }

func (t nodeQ2[E, B]) isLeaf() bool                    { return false }
func (t nodeQ2[E, B]) leftTree() evalTreeNode          { return t.left }
func (t nodeQ2[E, B]) rightTree() evalTreeNode         { return t.right }
func (t nodeQ2[E, B]) rightTree_() evalRightNode[E, B] { return t.right }

// A tree node where only the input type A is statically known.
type evalLeftNode[E any, A any] interface {
	evalTreeNode
}

// A tree node where only the output type B is statically known.
type evalRightNode[E any, B any] interface {
	evalTreeNode
	rightTree_() evalRightNode[E, B]

	qApply(start any) Eff[E, B]
	qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, B]
}

type evalQApply interface {
	apply(start any) effBase
}

type evalQ[E any, A any, B any] interface {
	evalLeftNode[E, A]
	evalRightNode[E, B]
}

func (_ leafQ[E, A, B]) evalQ(eff E, input A) B    { panic("marker method") }
func (_ nodeQ[E, A, X, B]) evalQ(eff E, input A) B { panic("marker method") }

func (l leafQ[E, A, B]) apply(start any) effBase {
	return l.lifted(start.(A))
}

func (l leafQ[E, A, B]) qApply(start any) Eff[E, B] {
	return l.lifted(start.(A))
}

func (l leafQ[E, A, B]) qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, B] {
	return newCont[E, Start, B](Union[E, Start](effect), concatQ[E, Start, A, B](queue.(evalQueue[E, Start, A]), l))
}

func (t nodeQ[E, A, X, B]) qApply(start any) Eff[E, B] {
	return qApply(start.(A), evalQueue[E, A, B](t))
}

func (t nodeQ[E, A, X, B]) qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, B] {
	return newCont[E, Start, B](Union[E, Start](effect), concatQ[E, Start, A, B](queue.(evalQueue[E, Start, A]), t))
}

func (t nodeQ2[E, B]) qApply(start any) Eff[E, B] {
	return qApply2[E, B](start, evalRightNode[E, B](t))
}

func (t nodeQ2[E, B]) qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, B] {
	return newCont[E, Start, B](Union[E, Start](effect), concatQ2Out[E, B](queue, t).(evalQueue[E, Start, B]))
}

type leafQ[E any, A any, B any] struct {
	lifted Arr[E, A, B]
}

type nodeQ[E any, A any, X any, B any] struct {
	left  evalQueue[E, A, X]
	right evalQueue[E, X, B]
}

type nodeQ2[E any, B any] struct {
	left  evalTreeNode
	right evalRightNode[E, B]
}

func liftQ[E any, A any, B any](f Arr[E, A, B]) evalQueue[E, A, B] {
	return leafQ[E, A, B]{lifted: f}
}

func concatQ[E any, A any, X any, B any](a2x evalQueue[E, A, X], x2b evalQueue[E, X, B]) evalQueue[E, A, B] {
	return nodeQ[E, A, X, B]{left: a2x, right: x2b}
}

func concatQ2Out[E any, B any](a2x evalTreeNode, x2b evalRightNode[E, B]) evalRightNode[E, B] {
	return nodeQ2[E, B]{left: a2x, right: x2b}
}

func qApply[E any, A any, B any](start A, q evalQueue[E, A, B]) Eff[E, B] {
	switch q := q.(type) {
	case leafQ[E, A, B]:
		return q.qApply(start)
	default:
		// q hasType nodeQ[E, A, X, B] where exists(X):
		return qApplyWithContinuation[E, A, B](
			start,          // A
			q.leftTree(),   // evalQueue[E, A, X]
			q.rightTree_(), // evalQueue[E, X, B]
		) // => Eff[E, B]
	}
}

func qApply2[E any, B any](start any, q evalRightNode[E, B]) Eff[E, B] {
	if q.isLeaf() {
		// start hasType A
		// q hasType leafQ[E, A, B]
		return q.qApply(start)
	} else {
		// q hasType nodeQ[E, A, X, B] where exists(X):
		return qApplyWithContinuation[E, any /* A */, B](
			start,          // A
			q.leftTree(),   // evalQueue[E, A, X]
			q.rightTree_(), // evalQueue[E, X, B]
		) // => Eff[E, B]
	}
}

func qApplyWithContinuation[E any, A any, B any](start A, tl evalLeftNode[E, A], tr evalRightNode[E, B]) Eff[E, B] {
	// (tl hasType evalQueue[E, A, X]
	//  tr hasType evalQueue[E, X, B]) where exists(X)

	if tl.isLeaf() {
		// tl hasType leafQ[E, A, X]
		// therefore: tl.lifted(start) hasType Eff[E, X]
		return qBind2[E, B]( // Eff[E, X] => evalQueue[E, X, B] => Eff[E, B]
			tl.(evalQApply).apply(start), // (start: A) => (tl.apply: A => Eff[E, X]) => Eff[E, X]
			tr,                           // evalQueue[E, X, B]
		) // => Eff[E, B]
	} else {
		// tl hasType nodeQ[E, A, Y, X] where exists(Y)
		// therefore tl.left  hasType evalQueue[E, A, Y]
		//           tl.right hasType evalQueue[E, Y, X]
		return qApplyWithContinuation[E, A, B](
			start,         // A
			tl.leftTree(), // evalQueue[E, A, Y]
			concatQ2Out[E, B](
				tl.rightTree(), // evalQueue[E, Y, X]
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

func qBind2[E any, B any](e effBase, k evalRightNode[E, B]) Eff[E, B] {
	// (e hasType Eff[E, X]
	//  k hasType evalQueue[E, X, B]) where exists(X)

	if m := e.impl(); m.isPure() {
		// m hasType Pure[E, X]
		// k hasType evalQueue[E, X, B]
		// k.qApply(...) hasType Eff[E, B]
		fmt.Printf("qBind2\n    %#v\n    %#v\n", e, k)
		return k.qApply(m.pureValue())
	} else {
		// m hasType Cont[E, Start, X]:
		// k hasType evalQueue[E, X, B]
		// k.qPrepend(...) hasType Eff[E, B]
		return k.qPrepend(m.getEffect(), m.getQueue())
	}
}
