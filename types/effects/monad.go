package effects

// Non-generic workaround.
type effBase interface {
	effBase()
	impl() effImplBase
}

type effImplBase interface {
	effImplBase()
	isPure() bool
	pureValue() any
	getEffect() EffectTag
	getQueue() evalTreeNode
}

func (_ Eff[_, _]) effBase()                   {}
func (e Eff[_, _]) impl() effImplBase          { return e.EffImpl }
func (_ Pure[_, _]) effImplBase()              {}
func (_ Cont[_, _, _]) effImplBase()           {}
func (_ Pure[_, _]) isPure() bool              { return true }
func (_ Cont[_, _, _]) isPure() bool           { return false }
func (p Pure[_, _]) pureValue() any            { return p.value }
func (_ Cont[_, _, _]) pureValue() any         { panic("not pure") }
func (_ Pure[_, _]) getEffect() EffectTag      { panic("pure") }
func (c Cont[_, _, _]) getEffect() EffectTag   { return EffectTag(c.effect) }
func (_ Pure[_, _]) getQueue() evalTreeNode    { panic("pure") }
func (c Cont[_, _, _]) getQueue() evalTreeNode { return evalTreeNode(c.queue) }

// Generic type.
type Eff[E any, T any] struct {
	EffImpl[E, T]
}

func asEff[E any, T any](impl EffImpl[E, T]) Eff[E, T] {
	return Eff[E, T]{impl}
}

type EffImpl[E any, T any] interface {
	eff(eff E) T // marker method - do not call
	effImplBase
}

// Marker type used for better compiler errors.
type effIsNotAnEffImpl struct{}

// Force different signature to prevent Eff from being used as an EffImpl,
// even though it embeds an EffImpl.
func (_ Eff[E, T]) eff(_ effIsNotAnEffImpl) {}

func (_ Pure[E, T]) eff(eff E) T    { panic("marker method") }
func (_ Cont[E, _, T]) eff(eff E) T { panic("marker method") }

type Pure[E any, T any] struct {
	value T
}

func newPure[E any, T any](value T) Eff[E, T] {
	return asEff[E, T](Pure[E, T]{value: value})
}

func RunPure[E any, T any](e Eff[E, T]) (T, error) {
	panic("not implemented")
}

func RunPureOrFail[E any, T any](e Eff[E, T]) T {
	switch m := e.EffImpl.(type) {
	case Pure[E, T]:
		return m.value
	default:
		panic("unhandled effect")
	}
}

type Start any

type Union[E any, T any] EffectTag

func inj[E any, T any, L EffectTag](tag L) Union[E, T] {
	return Union[E, T](tag)
}

type Cont[E any, A any, B any] struct {
	effect Union[E, A]
	queue  evalQueue[E, A, B]
}

func newCont[E any, A any, B any](effect Union[E, A], queue evalQueue[E, A, B]) Eff[E, B] {
	return asEff[E, B](Cont[E, A, B]{
		effect: effect,
		queue:  queue,
	})
}

func injectEffect[E any, T any, L EffectTag](tag L) Eff[E, T] {
	// We use Start/any as a workaround for not having universal quantification
	return newCont[E, Start, T](
		inj[E, Start](tag),
		liftQ(func(value Start) Eff[E, T] { return newPure[E](value.(T)) }),
	)
}

func Lift[E any, A any, B any](f func(arg A) B) func(arg Eff[E, A]) Eff[E, B] {
	return func(arg Eff[E, A]) Eff[E, B] {
		return Return[E](f(arg.EffImpl.(Pure[E, A]).value))
	}
}

func Return[E any, T any](value T) Eff[E, T] {
	return newPure[E](value)
}

func Map[E any, A any, B any](arg Eff[E, A], f func(arg A) B) Eff[E, B] {
	switch m := arg.EffImpl.(type) {
	case Pure[E, A]:
		return newPure[E](f(m.value))
	case Cont[E, Start, A]: // hack for missing universal quantification
		g := func(arg A) Eff[E, B] {
			return newPure[E](f(arg))
		}
		return newCont(m.effect, concatQ(m.queue, liftQ(g)))
	default:
		panic("unreachable")
	}
}

func FlatMap[E any, A any, B any](arg Eff[E, A], f func(arg A) Eff[E, B]) Eff[E, B] {
	switch m := arg.EffImpl.(type) {
	case Pure[E, A]:
		return f(m.value)
	case Cont[E, Start, A]: // hack for missing universal quantification
		return newCont(m.effect, concatQ(m.queue, liftQ(f)))
	default:
		panic("unreachable")
	}
}

// Appends the effects from `others`, but keeps the value from `first`.
func Chain[E any, A any, Discard any](first Eff[E, A], others ...Eff[E, Discard]) Eff[E, A] {
	r := first
	for _, other := range others {
		r = r.And(other.Discard())
	}
	return r
}

// Discards the value but keeps the effects.
func (e Eff[E, A]) Discard() Eff[E, Unit] {
	return Map(e, func(_ A) Unit { return Unit{} })
}

// Appends the effects from other, but keeps the current value from e.
func (e Eff[E, A]) And(other Eff[E, Unit]) Eff[E, A] {
	return FlatMap(e, func(arg A) Eff[E, A] {
		return Map(other, func(_ Unit) A { return arg })
	})
}

// Appends the effects returned by f, but keeps the current value from e.
func (e Eff[E, A]) Then(f func(arg A) Eff[E, Unit]) Eff[E, A] {
	return FlatMap(e, func(arg A) Eff[E, A] {
		return Map(f(arg), func(_ Unit) A { return arg })
	})
}
