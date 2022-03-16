package effects 

type Arr[E any, A any, B any] func(arg A) Eff[E, B]

// evalQueue represents a type-aligned sequence of Kleisli arrows.
type evalQueue[E any, A any, B any] interface {
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
	leftTree_() evalLeftNode[E, A]
}

// A tree node where only the output type B is statically known.
type evalRightNode[E any, B any] interface {
	evalTreeNode
	rightTree_() evalRightNode[E, B]

	qApply(input any) Eff[E, B]
	qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, B]
}

func (l identQ[E, A]) apply(input any) effBase         { return l.qApply(input) }
func (l identQ[E, A]) isLeaf() bool                    { return true }
func (l identQ[E, A]) leftTree() evalTreeNode          { return nil }
func (l identQ[E, A]) leftTree_() evalLeftNode[E, A]   { return nil }
func (l identQ[E, A]) rightTree() evalTreeNode         { return nil }
func (l identQ[E, A]) rightTree_() evalRightNode[E, A] { return nil }

func (l leafQ[E, A, B]) apply(input any) effBase         { return l.qApply(input) }
func (l leafQ[E, A, B]) isLeaf() bool                    { return true }
func (l leafQ[E, A, B]) leftTree() evalTreeNode          { return nil }
func (l leafQ[E, A, B]) leftTree_() evalLeftNode[E, A]   { return nil }
func (l leafQ[E, A, B]) rightTree() evalTreeNode         { return nil }
func (l leafQ[E, A, B]) rightTree_() evalRightNode[E, B] { return nil }

func (t nodeQErased[E, B]) isLeaf() bool                    { return false }
func (t nodeQErased[E, B]) leftTree() evalTreeNode          { return t.left }
func (t nodeQErased[E, B]) rightTree() evalTreeNode         { return t.right }
func (t nodeQErased[E, B]) rightTree_() evalRightNode[E, B] { return t.right }

type identQ[E any, A any] struct{}

type leafQ[E any, A any, B any] struct {
	lifted   Arr[E, A, B]
	debugTag string
}

type nodeQErased[E any, B any] struct {
	left  evalTreeNode
	right evalRightNode[E, B]
}

func passThruQ[E any, A any]() evalQueue[E, A, A] {
	return identQ[E, A]{}
	// return liftQ(newPure[E, A], "PassThru")
}

func liftQ[E any, A any, B any](f Arr[E, A, B], debugTag string) evalQueue[E, A, B] {
	return leafQ[E, A, B]{lifted: f, debugTag: debugTag}
}

func composeRunQ[E any, B any, C any](a2b evalRightNode[E, B], b2c func(eff Eff[E, B]) Eff[E, C], debugTag string) evalRightNode[E, C] {
	return liftQ(qCompose[E, any](a2b, b2c), debugTag)
}

func composeQ[E any, X any, B any](a2x evalRightNode[E, X], x2b evalQueue[E, X, B]) evalRightNode[E, B] {
	return nodeQErased[E, B]{left: a2x, right: x2b}
}

// composeQErased is a partially type-erased version of composeQ.
func composeQErased[E any, B any](a2x evalTreeNode, x2b evalRightNode[E, B]) evalRightNode[E, B] {
	return nodeQErased[E, B]{left: a2x, right: x2b}
}

func (l identQ[E, A]) qApply(start any /* A */) Eff[E, A] {
	return newPure[E, A](start.(A))
}

func (l leafQ[E, A, B]) qApply(start any /* A */) Eff[E, B] {
	return l.lifted(start.(A))
}

func (t nodeQErased[E, B]) qApply(start any /* X */) Eff[E, B] {
	return qApplyWithContinuation[E, B](
		start,   // X
		t.left,  // evalQueue[E, X, Y]
		t.right, // evalQueue[E, Y, B]
	) // => Eff[E, B]
}

func qApplyWithContinuation[E any, B any](start any, tl evalTreeNode, tr evalRightNode[E, B]) Eff[E, B] {
	// (tl hasType evalQueue[E, X, Z]
	//  tr hasType evalQueue[E, Z, B]) where exists(X) & exists(Z)

	head, tail := qExtractHeadTail(tl, tr)
	return qBindErased[E, B](head.apply(start), tail)
}

func qExtractHeadTail[E any, B any](tl evalTreeNode /* X => Y */, tr evalRightNode[E, B] /* Y => B */) (head evalQApply, tail evalRightNode[E, B]) {
	// (tl hasType evalQueue[E, X, Z]
	//  tr hasType evalQueue[E, Z, B]) where exists(X) & exists(Z)

	if tl.isLeaf() {
		return tl.(evalQApply), tr
	} else {
		return qExtractHeadTail[E, B](
			tl.leftTree(),                            // evalQueue[E, X, Y]
			composeQErased[E, B](tl.rightTree(), tr), // evalQueue[E, Y, Z] => evalQueue[E, Z, B] => evalQueue[E, X, B]
		)
	}
}

func qCompose[E any, A any, B any, C any](a2b evalRightNode[E, B], b2c func(eff Eff[E, B]) Eff[E, C]) Arr[E, A, C] {
	return func(input A) Eff[E, C] {
		return b2c(a2b.qApply(input)) // => Eff[E, C]
	} // => (A => Eff[E, C])
}

func qBindErased[E any, B any](e effBase, k evalRightNode[E, B]) Eff[E, B] {
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

func (l identQ[E, A]) qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, A] {
	return newContUnchecked(effect, composeQErased[E, A](queue, l))
}

func (l leafQ[E, A, B]) qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, B] {
	return newContUnchecked(effect, composeQErased[E, B](queue, l))
}

func (t nodeQErased[E, B]) qPrepend(effect EffectTag, queue evalTreeNode) Eff[E, B] {
	return newContUnchecked(effect, composeQErased[E, B](queue, t))
}
