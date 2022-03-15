package effects 

type Arr[E any, A any, B any] func(arg A) Eff[E, B]

// evalQueue represents a type-aligned sequence of Kleisli arrows.
type evalQueue[E any, A any, B any] interface {
	evalQ(effects E, input A) B // marker method - do not call
	applyTo(input A) Eff[E, B]

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

type evalQApply interface {
	apply(input any) effBase
}

// A tree node where only the input type A is statically known.
type evalLeftNode[E any, A any] interface {
	evalTreeNode
}

// A tree node where only the output type B is statically known.
type evalRightNode[E any, B any] interface {
	evalTreeNode
	rightTree_() evalRightNode[E, B]

	qApply(input any) Eff[E, B]
	qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, B]
}

func (l leafQ[E, A, B]) evalQ(effects E, input A) B      { panic("marker method") }
func (l leafQ[E, A, B]) apply(input any) effBase         { return l.qApply(input) }
func (l leafQ[E, A, B]) applyTo(input A) Eff[E, B]       { return l.qApply(input) }
func (l leafQ[E, A, B]) isLeaf() bool                    { return true }
func (l leafQ[E, A, B]) leftTree() evalTreeNode          { return nil }
func (l leafQ[E, A, B]) leftTree_() evalLeftNode[E, B]   { return nil }
func (l leafQ[E, A, B]) rightTree() evalTreeNode         { return nil }
func (l leafQ[E, A, B]) rightTree_() evalRightNode[E, B] { return nil }

func (t nodeQ[E, A, X, B]) evalQ(effects E, input A) B      { panic("marker method") }
func (t nodeQ[E, A, X, B]) applyTo(input A) Eff[E, B]       { return t.qApply(input) }
func (t nodeQ[E, A, X, B]) isLeaf() bool                    { return false }
func (t nodeQ[E, A, X, B]) leftTree() evalTreeNode          { return t.left }
func (t nodeQ[E, A, X, B]) leftTree_() evalLeftNode[E, A]   { return t.left }
func (t nodeQ[E, A, X, B]) rightTree() evalTreeNode         { return t.right }
func (t nodeQ[E, A, X, B]) rightTree_() evalRightNode[E, B] { return t.right }

func (t nodeQErased[E, B]) isLeaf() bool                    { return false }
func (t nodeQErased[E, B]) leftTree() evalTreeNode          { return t.left }
func (t nodeQErased[E, B]) rightTree() evalTreeNode         { return t.right }
func (t nodeQErased[E, B]) rightTree_() evalRightNode[E, B] { return t.right }

type leafQ[E any, A any, B any] struct {
	lifted Arr[E, A, B]
}

type nodeQ[E any, A any, X any, B any] struct {
	left  evalQueue[E, A, X]
	right evalQueue[E, X, B]
}

// nodeQErased is a partially type-erased version of nodeQ
// that is only used during the evaluation process.
type nodeQErased[E any, B any] struct {
	left  evalTreeNode
	right evalRightNode[E, B]
}

func liftQ[E any, A any, B any](f Arr[E, A, B]) evalQueue[E, A, B] {
	return leafQ[E, A, B]{lifted: f}
}

func composeQ[E any, A any, X any, B any](a2x evalQueue[E, A, X], x2b evalQueue[E, X, B]) evalQueue[E, A, B] {
	return nodeQ[E, A, X, B]{left: a2x, right: x2b}
}

// composeQErased is a partially type-erased version of composeQ.
func composeQErased[E any, B any](a2x evalTreeNode, x2b evalRightNode[E, B]) evalRightNode[E, B] {
	return nodeQErased[E, B]{left: a2x, right: x2b}
}

func (l leafQ[E, A, B]) qApply(start any) Eff[E, B] {
	return l.lifted(start.(A))
}

func (t nodeQ[E, A, X, B]) qApply(start any) Eff[E, B] {
	return qApplyWithContinuation[E, A, B](
		start.(A), // A
		t.left,    // evalQueue[E, A, X]
		t.right,   // evalQueue[E, X, B]
	) // => Eff[E, B]
}

func (t nodeQErased[E, B]) qApply(start any) Eff[E, B] {
	return qApplyWithContinuation[E, any /* A */, B](
		start,   // A
		t.left,  // evalQueue[E, A, X]
		t.right, // evalQueue[E, X, B]
	) // => Eff[E, B]
}

func qApplyWithContinuation[E any, A any, B any](start A, tl evalLeftNode[E, A], tr evalRightNode[E, B]) Eff[E, B] {
	// (tl hasType evalQueue[E, A, X]
	//  tr hasType evalQueue[E, X, B]) where exists(X)

	head, tail := qExtractHeadTail(tl, tr)
	return qBind2[E, B](head.apply(start), tail)
}

func qExtractHeadTail[E any, A any, B any](tl evalLeftNode[E, A], tr evalRightNode[E, B]) (head evalQApply, tail evalRightNode[E, B]) {
	// (tl hasType evalQueue[E, A, X]
	//  tr hasType evalQueue[E, X, B]) where exists(X)

	if tl.isLeaf() {
		return tl.(evalQApply), tr
	} else {
		return qExtractHeadTail[E, A, B](
			tl.leftTree(),                            // evalQueue[E, A, Y]
			composeQErased[E, B](tl.rightTree(), tr), // evalQueue[E, Y, X] => evalQueue[E, X, B] => evalQueue[E, Y, B]
		)
	}
}

func qCompose[E any, A any, B any, C any](a2b evalQueue[E, A, B], b2c func(eff Eff[E, B]) Eff[E, C]) Arr[E, A, C] {
	return func(input A) Eff[E, C] {
		return b2c(a2b.applyTo(input)) // => Eff[E, C]
	} // => (A => Eff[E, C])
}

func qBind[E any, A any, B any](e Eff[E, A], k evalQueue[E, A, B]) Eff[E, B] {
	switch m := e.EffImpl.(type) {
	case Pure[E, A]:
		return k.applyTo(m.value)
	case Cont[E, Start, A]:
		return newCont(m.effect, composeQ(m.queue, k))
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
		return k.qApply(m.pureValue())
	} else {
		// m hasType Cont[E, Start, X]:
		// k hasType evalQueue[E, X, B]
		// k.qPrepend(...) hasType Eff[E, B]
		return k.qPrepend(m.getEffect(), m.getQueue())
	}
}

func (l leafQ[E, A, B]) qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, B] {
	return newCont[E, Start, B](Union[E, Start](effect), composeQ[E, Start, A, B](queue.(evalQueue[E, Start, A]), l))
}

func (t nodeQ[E, A, X, B]) qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, B] {
	return newCont[E, Start, B](Union[E, Start](effect), composeQ[E, Start, A, B](queue.(evalQueue[E, Start, A]), t))
}

func (t nodeQErased[E, B]) qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, B] {
	panic("not implemented")
	// return newCont[E, Start, B](Union[E, Start](effect), composeQErased[E, B](queue, t).(evalQueue[E, Start, B]))
}
