package effects

import (
	"github.com/cr7pt0gr4ph7/functional-go/types/list"
)

// Represents the empty tuple.
type Unit struct{}

var UnitValue Unit

// Constraint type for effects markers.
type Effect interface {
	effect()
}

func (_ ReaderI[E, R]) effect()       {}
func (_ WriterI[E, W]) effect()       {}
func (_ StateI[E, S]) effect()        {}
func (_ CoroutineI[E, Y, R]) effect() {}

// Interface type for effect representations.
type EffectTag interface {
	effectTag()
}

func (_ AskEffect[R]) effectTag()      {}
func (_ TellEffect[W]) effectTag()     {}
func (_ GetEffect[S]) effectTag()      {}
func (_ SetEffect[S]) effectTag()      {}
func (_ YieldEffect[Y, R]) effectTag() {}

// Interface type for reprsentations of effects that result in a value of type T.
type TypedEffectTag[T any] interface {
	EffectTag
	effectResult() T // marker method - do not call
}

func ApplyContinuationToEffectResult[L TypedEffectTag[A], E any, A any, B any](effect L, continuation evalRightNode[E, B], effectResult A) Eff[E, B] {
	return continuation.qApply(effectResult)
}

// ===================
// :: Reader Effect ::
// ===================

// Effect: Read a shared immutable value from the environment.
type Reader[E any, R any] interface {
	Effect
	Ask() Eff[E, R]
}

func _[E Reader[E, R], R any]() {
	// Statically ensure that certain interfaces are implemented correctly
	var _ Reader[E, R] = ReaderI[E, R]{}
	var _ TypedEffectTag[R] = AskEffect[R]{}
}

// DSL implementation for `Reader[E, R]`.
type ReaderI[E Reader[E, R], R any] struct{}

func (_ ReaderI[E, R]) Ask() Eff[E, R] {
	return injectEffect[E, R](AskEffect[R]{})
}

// Effect tag for `Reader[E, R].Ask() Eff[E, R]`.
type AskEffect[R any] struct{}

func (_ AskEffect[R]) effectResult() R { panic("marker method") }

func RunReader[R any, E Reader[E, R]](value R, e Eff[E, R]) Eff[E, R] {
	loop := func(e Eff[E, R]) Eff[E, R] {
		return RunReader(value, e)
	}

	switch m := e.EffImpl.(type) {
	case Cont[E, R]:
		switch t := m.effect.(type) {
		case AskEffect[R]:
			return RunReader(value, ApplyContinuationToEffectResult(t, m.queue, value))
		default:
			return newContUnchecked(t, composeRunQ(m.queue, loop, "RunReader"))
		}
	}
	return e
}

// ===================
// :: Writer Effect ::
// ===================

// Effect: Send outputs to the effects environment.
type Writer[E any, W any] interface {
	Effect
	Tell(output W) Eff[E, Unit]
}

func _[E Writer[E, W], W any]() {
	// Statically ensure that certain interfaces are implemented correctly
	var _ Writer[E, W] = WriterI[E, W]{}
	var _ TypedEffectTag[Unit] = TellEffect[W]{}
}

// DSL implementation for `Writer[E, W]`.
type WriterI[E Writer[E, W], W any] struct{}

func (_ WriterI[E, W]) Tell(output W) Eff[E, Unit] {
	return injectEffect[E, Unit](TellEffect[W]{output: output})
}

// Effect tag for `Writer[E, W].Tell(output W) Eff[E, Unit]`.
type TellEffect[W any] struct{ output W }

func (_ TellEffect[W]) effectResult() Unit { panic("marker method") }

type WriterResult[T any, W any] struct {
	Value   T
	Written list.List[W]
}

func RunWriter[W any, E Writer[E, W], T any](e Eff[E, T]) Eff[E, WriterResult[T, W]] {
	switch m := e.EffImpl.(type) {
	case Pure[E, T]:
		return newPure[E](WriterResult[T, W]{
			Value:   m.value,
			Written: list.Empty[W](),
		})
	case Cont[E, T]:
		switch t := m.effect.(type) {
		case TellEffect[W]:
			kx := RunWriter[W](ApplyContinuationToEffectResult(t, m.queue, UnitValue))
			return FlatMap(kx, func(x WriterResult[T, W]) Eff[E, WriterResult[T, W]] {
				return newPure[E](WriterResult[T, W]{
					Value:   x.Value,
					Written: x.Written.Push(t.output),
				})
			})
		default:
			return newContUnchecked(t, composeRunQ(m.queue, RunWriter[W, E, T], "RunWriter"))
		}
	default:
		panic("unreachable")
	}
}

func RunWriterReverse[W any, E Writer[E, W], T any](l List[W], e Eff[E, T]) Eff[E, WriterResult[T, W]] {
	switch m := e.EffImpl.(type) {
	case Pure[E, T]:
		return newPure[E](WriterResult[T, W]{
			Value:   m.value,
			Written: l,
		})
	case Cont[E, T]:
		switch t := m.effect.(type) {
		case TellEffect[W]:
			return RunWriterReverse[W](l.Push(t.output), ApplyContinuationToEffectResult(t, m.queue, UnitValue))
		default:
			loop := func(e Eff[E, T]) Eff[E, WriterResult[T, W]] {
				return RunWriterReverse(l, e)
			}
			return newContUnchecked(t, composeRunQ(m.queue, loop, "RunWriterReverse"))
		}
	default:
		panic("unreachable")
	}
}

// ==================
// :: State Effect ::
// ==================

// Effect: Provides read/write access to a shared updatable state value of type S.
type State[E any, S any] interface {
	Effect
	Get() Eff[E, S]
	Set(newState S) Eff[E, Unit]
}

func _[E State[E, S], S any]() {
	// Statically ensure that certain interfaces are implemented correctly
	var _ State[E, S] = StateI[E, S]{}
	var _ TypedEffectTag[S] = GetEffect[S]{}
	var _ TypedEffectTag[Unit] = SetEffect[S]{}
}

// DSL implementation for `State[E, S]`.
type StateI[E State[E, S], S any] struct{}

func (_ StateI[E, S]) Get() Eff[E, S] {
	return injectEffect[E, S](GetEffect[S]{})
}

func (_ StateI[E, S]) Set(newState S) Eff[E, Unit] {
	return injectEffect[E, Unit](SetEffect[S]{newState: newState})
}

// Effect tag for `State[E, S].Get() Eff[E, S]`.
type GetEffect[S any] struct{}

func (_ GetEffect[S]) effectResult() S { panic("marker method") }

// Effect tag for `State[_, S].Set(newState S) Eff[E, Unit]`.
type SetEffect[S any] struct{ newState S }

func (_ SetEffect[S]) effectResult() Unit { panic("marker method") }

type StateResult[T any, S any] struct {
	Value T
	State S
}

func RunState[S any, E State[E, S], T any](state S, e Eff[E, T]) Eff[E, StateResult[T, S]] {
	fmt.Println("RunState", state, e)

	switch m := e.EffImpl.(type) {
	case Pure[E, T]:
		return newPure[E](StateResult[T, S]{
			Value: m.value,
			State: state,
		})
	case Cont[E, T]:
		switch t := m.effect.(type) {
		case GetEffect[S]:
			return RunState(state, ApplyContinuationToEffectResult(t, m.queue, state))
		case SetEffect[S]:
			return RunState(t.newState, ApplyContinuationToEffectResult(t, m.queue, UnitValue))
		default:
			loop := func(e Eff[E, T]) Eff[E, StateResult[T, S]] {
				return RunState(state, e)
			}
			return newContUnchecked(t, composeRunQ(m.queue, loop, "RunState"))
		}
	default:
		panic("unreachable")
	}
}

// ======================
// :: Coroutine Effect ::
// ======================

// Effect: A type representing a yielding of control.
//
// The type variables have the following meaning:
//
// * A: The current type.
// * Y: The input to the continuation function.
// * R: The output of the continuation.
type Coroutine[E any, Y any, R any] interface {
	Effect
	Yield(output Y) Eff[E, R]
}

func _[E Coroutine[E, Y, R], Y any, R any]() {
	// Statically ensure that certain interfaces are implemented correctly
	var _ Coroutine[E, Y, R] = CoroutineI[E, Y, R]{}
	var _ TypedEffectTag[R] = YieldEffect[Y, R]{}
}

// DSL implementation for `Coroutine[E, Y, R]`.
type CoroutineI[E Coroutine[E, Y, R], Y any, R any] struct{}

func (_ CoroutineI[E, Y, R]) Yield(output Y) Eff[E, R] {
	return injectEffect[E, R](YieldEffect[Y, R]{output: output})
}

// Effect tag for `Couroutine[E, Y, R].Yield(output Y) Eff[E, R]`.
type YieldEffect[Y any, R any] struct{ output Y }

func (_ YieldEffect[Y, R]) effectResult() R { panic("marker method") }
