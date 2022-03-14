package effects

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

func (_ AskEffect[E, R]) effectTag()  {}
func (_ TellEffect[E, W]) effectTag() {}
func (_ GetEffect[E, W]) effectTag()  {}
func (_ SetEffect[E, W]) effectTag()  {}

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
	return injectEffect[E, R](AskEffect[E, R]{})
}

type AskEffect[E any, R any] struct{}

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
	return injectEffect[E, Unit](TellEffect[E, W]{output: output})
}

type TellEffect[E any, W any] struct{ output W }

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
	return injectEffect[E, S](GetEffect[E, S]{})
}

func (_ StateI[E, S]) Set(newState S) Eff[E, Unit] {
	return injectEffect[E, Unit](SetEffect[E, S]{newState: newState})
}

type GetEffect[E any, S any] struct{}

type SetEffect[E any, S any] struct{ newState S }
