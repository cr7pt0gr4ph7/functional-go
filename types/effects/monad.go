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

func (_ Eff[_, _]) effBase()          {}
func (e Eff[_, _]) impl() effImplBase { return e.EffImpl }

func (_ Pure[_, _]) effImplBase()           {}
func (_ Pure[_, _]) isPure() bool           { return true }
func (p Pure[_, _]) pureValue() any         { return p.value }
func (_ Pure[_, _]) getEffect() EffectTag   { panic("pure") }
func (_ Pure[_, _]) getQueue() evalTreeNode { panic("pure") }

func (_ Cont[_, _]) effImplBase()           {}
func (_ Cont[_, _]) isPure() bool           { return false }
func (_ Cont[_, _]) pureValue() any         { panic("not pure") }
func (c Cont[_, _]) getEffect() EffectTag   { return EffectTag(c.effect) }
func (c Cont[_, _]) getQueue() evalTreeNode { return evalTreeNode(c.queue) }

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

func (_ Pure[E, T]) eff(eff E) T { panic("marker method") }
func (_ Cont[E, T]) eff(eff E) T { panic("marker method") }

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

type ValueFromEffect any

type Cont[E any, B any] struct {
	effect EffectTag
	queue  evalQueue[E, ValueFromEffect, B]
}

func newCont[E any, B any](effect EffectTag, queue evalQueue[E, ValueFromEffect, B]) Eff[E, B] {
	return asEff[E, B](Cont[E, B]{
		effect: effect,
		queue:  queue,
	})
}

func newPureFromEffect[E any, T any](value ValueFromEffect) Eff[E, T] {
	return newPure[E, T](value.(T))
}

func injectEffect[E any, T any, L EffectTag](tag L) Eff[E, T] {
	return newCont[E, T](tag, liftQ(newPureFromEffect[E, T]))
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
	case Cont[E, A]:
		g := func(arg A) Eff[E, B] {
			return newPure[E](f(arg))
		}
		return newCont(m.effect, composeQ(m.queue, liftQ(g)))
	default:
		if !m.isPure() {

		}
		panic("unreachable")
	}
}

func FlatMap[E any, A any, B any](arg Eff[E, A], f func(arg A) Eff[E, B]) Eff[E, B] {
	switch m := arg.EffImpl.(type) {
	case Pure[E, A]:
		return f(m.value)
	case Cont[E, A]:
		return newCont(m.effect, composeQ(m.queue, liftQ(f)))
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
