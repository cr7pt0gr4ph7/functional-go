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

func (_ ReaderE[E, R]) effectTag() {}
func (_ WriterE[E, W]) effectTag() {}
func (_ StateE[E, W]) effectTag()  {}

// Effect: Read a shared immutable value from the environment.
type Reader[E any, R any] interface {
	Effect
	Ask() Eff[E, R]
}

func _[E Reader[E, R], R any]() {
	// Ensure that ReaderI[E,R] implements Reader[E,R]
	var _ Reader[E, R] = ReaderI[E, R]{}
}

type ReaderE[E any, R any] struct{}

type ReaderI[E Reader[E, R], R any] struct{}

func (_ ReaderI[E, R]) Ask() Eff[E, R] {
	panic("not implemented")
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

type WriterE[E any, W any] struct{}

type WriterI[E Writer[E, W], W any] struct{}

func (_ WriterI[E, W]) Tell(output W) Eff[E, Unit] {
	panic("not implemented")
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

type StateE[E any, S any] struct{}

type StateI[E State[E, S], S any] struct{}

func (_ StateI[E, S]) Get() Eff[E, S] {
	panic("not implemented")
}

func (_ StateI[E, S]) Set(newState S) Eff[E, Unit] {
	panic("not implemented")
}
