package effects

// Represents the empty tuple.
type Unit struct{}

// Constraint type for effects markers.
type Effect interface {
	effect()
}

type Apply[E Effect, T any] struct {
}

// Effect: Read a shared immutable value from the environment.
type Reader[R any] interface {
	Effect
	Ask() Apply[Reader[R], R]
}

// Effect: Send outputs to the effects environment.
type Writer[W any] interface {
	Effect
	Tell(output W) Apply[Writer[W], Unit]
}

// Effect: Provides read/write access to a shared updatable state value of type S.
type State[S any] interface {
	Effect
	Get() Apply[State[S], S]
	Set(newState S) Apply[State[S], Unit]
}
