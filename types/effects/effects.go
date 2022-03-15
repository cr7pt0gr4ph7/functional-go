package effects

import (
	"github.com/cr7pt0gr4ph7/functional-go/types/list"
)

// Represents the empty tuple.
type Unit struct{}

// Constraint type for effects markers.
type Effect interface {
	effect()
}

func (_ ReaderI[E, R]) effect() {}
func (_ WriterI[E, W]) effect() {}
func (_ StateI[E, W]) effect()  {}

// Interface type for effect representations.
type EffectTag interface {
	effectTag()
}

func (_ AskEffect[R]) effectTag()  {}
func (_ TellEffect[W]) effectTag() {}
func (_ GetEffect[S]) effectTag()  {}
func (_ SetEffect[S]) effectTag()  {}

// Effect: Read a shared immutable value from the environment.
type Reader[E any, R any] interface {
	Effect
	Ask() Eff[E, R]
}

func _[E Reader[E, R], R any]() {
	// Ensure that ReaderI[E,R] implements Reader[E,R]
	var _ Reader[E, R] = ReaderI[E, R]{}
}

type ReaderI[E Reader[E, R], R any] struct{}

func (_ ReaderI[E, R]) Ask() Eff[E, R] {
	return injectEffect[E, R](AskEffect[R]{})
}

type AskEffect[R any] struct{}

func RunReader[R any, E Reader[E, R]](value R, e Eff[E, R]) Eff[E, R] {
	loop := func(e Eff[E, R]) Eff[E, R] {
		return RunReader(value, e)
	}

	switch m := e.EffImpl.(type) {
	case Cont[E, Start, R]:
		switch t := m.effect.(type) {
		case AskEffect[R]:
			return RunReader(value, m.queue.qApply(Start(value)))
		default:
			return newCont(Union[E, Start](t), liftQ(qCompose(m.queue, loop)))
		}
	}
	return e
}

// Effect: Send outputs to the effects environment.
type Writer[E any, W any] interface {
	Effect
	Tell(output W) Eff[E, Unit]
}

func _[E Writer[E, W], W any]() {
	// Ensure that WriterI[E,W] implements Writer[E,W]
	var _ Writer[E, W] = WriterI[E, W]{}
}

type WriterI[E Writer[E, W], W any] struct{}

func (_ WriterI[E, W]) Tell(output W) Eff[E, Unit] {
	return injectEffect[E, Unit](TellEffect[W]{output: output})
}

type TellEffect[W any] struct{ output W }

type WriterResult[T any, W any] struct {
	Value   T
	Written list.List[W]
}

func RunWriter[W any, E Writer[E, W], T any](e Eff[E, T]) Eff[E, WriterResult[T, W]] {
	switch m := e.EffImpl.(type) {
	case Pure[E, T]:
		return newPure[E](WriterResult[T, W]{
			Value:   m.value,
			Written: list.Nil[W]{},
		})
	case Cont[E, Start, T]:
		k := qCompose(m.queue, RunWriter[W, E, T])
		switch t := m.effect.(type) {
		case TellEffect[W]:
			kx := k(Start(Unit{}))
			return FlatMap(kx, func(x WriterResult[T, W]) Eff[E, WriterResult[T, W]] {
				return newPure[E](WriterResult[T, W]{
					Value:   x.Value,
					Written: x.Written.Push(t.output),
				})
			})
		default:
			return newCont(Union[E, Start](t), liftQ(k))
		}
	default:
		panic("unreachable")
	}
}

// Effect: Provides read/write access to a shared updatable state value of type S.
type State[E any, S any] interface {
	Effect
	Get() Eff[E, S]
	Set(newState S) Eff[E, Unit]
}

func _[E State[E, S], S any]() {
	// Ensure that StateI[E,S] implements State[E,S]
	var _ State[E, S] = StateI[E, S]{}
}

type StateI[E State[E, S], S any] struct{}

func (_ StateI[E, S]) Get() Eff[E, S] {
	return injectEffect[E, S](GetEffect[S]{})
}

func (_ StateI[E, S]) Set(newState S) Eff[E, Unit] {
	return injectEffect[E, Unit](SetEffect[S]{newState: newState})
}

type GetEffect[S any] struct{}

type SetEffect[S any] struct{ newState S }
