package effects

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
	log.OnRunEffect("RunReader", value, e)

	switch m := e.EffImpl.(type) {
	case Cont[E, R]:
		switch t := m.effect.(type) {
		case AskEffect[R]:
			return RunReader(value, ApplyContinuationToEffectResult(t, m.queue, value))
		default:
			return ForwardEffectWithState(m, RunReader[R, E], value, "RunReader")
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

func _[E Writer[E, W], WL listBuilder[WL, W], W any, T any]() {
	// Statically ensure that certain interfaces are implemented correctly
	var _ Writer[E, W] = WriterI[E, W]{}
	var _ TypedEffectTag[Unit] = TellEffect[W]{}
	var _ Interpreter[E, T, WriterResult[T, WL]] = runWriter[WL, W, E, T]{}
	var _ Interpreter[E, T, WriterResult[T, WL]] = runWriterReverse[WL, W, E, T]{}
	var _ InterpreterWithState[runWriterReverse[WL, W, E, T], WL] = runWriterReverse[WL, W, E, T]{}
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
	Written W
}

type listBuilder[Self any, T any] interface {
	Push(item T) Self
}

func RunWriter[WL listBuilder[WL, W], W any, E Writer[E, W], T any](e Eff[E, T]) Eff[E, WriterResult[T, WL]] {
	return runWriter[WL, W, E, T]{}.Run(e)
}

type runWriter[WL listBuilder[WL, W], W any, E Writer[E, W], T any] struct{}

func (r runWriter[WL, W, E, T]) Name() string {
	return "RunWriter"
}

func (r runWriter[WL, W, E, T]) Run(e Eff[E, T]) Eff[E, WriterResult[T, WL]] {
	return RunImpl[E, T, WriterResult[T, WL]](r, e)
}

func (r runWriter[WL, W, E, T]) HandlePure(value T) WriterResult[T, WL] {
	// We expect the default value of WL to be a valid empty list instance
	var emptyList WL
	return WriterResult[T, WL]{
		Value:   value,
		Written: emptyList,
	}
}

func (r runWriter[WL, W, E, T]) HandleEffect(effect EffectTag, m Cont[E, T]) (_ Eff[E, WriterResult[T, WL]]) {
	switch t := m.effect.(type) {
	case TellEffect[W]:
		k := r.Run(ApplyContinuationToEffectResult(t, m.queue, UnitValue))

		return FlatMap(k, func(x WriterResult[T, WL]) Eff[E, WriterResult[T, WL]] {
			return newPure[E](WriterResult[T, WL]{
				Value:   x.Value,
				Written: x.Written.Push(t.output),
			})
		})
	}
	return
}

func RunWriterReverse[WL listBuilder[WL, W], W any, E Writer[E, W], T any](written WL, e Eff[E, T]) Eff[E, WriterResult[T, WL]] {
	return runWriterReverse[WL, W, E, T]{}.WithState(written).Run(e)
}

type runWriterReverse[WL listBuilder[WL, W], W any, E Writer[E, W], T any] struct {
	InterpreterWithStateImpl[WL, runWriterReverse[WL, W, E, T], *runWriterReverse[WL, W, E, T]]
}

func (r runWriterReverse[WL, W, E, T]) Name() string {
	return "RunWriterReverse"
}

func (r runWriterReverse[WL, W, E, T]) Run(e Eff[E, T]) Eff[E, WriterResult[T, WL]] {
	return RunImpl[E, T, WriterResult[T, WL]](r, e)
}

func (r runWriterReverse[WL, W, E, T]) HandlePure(value T) WriterResult[T, WL] {
	return WriterResult[T, WL]{
		Value:   value,
		Written: r.State(),
	}
}

func (r runWriterReverse[WL, W, E, T]) HandleEffect(effect EffectTag, m Cont[E, T]) (_ Eff[E, WriterResult[T, WL]]) {
	switch t := m.effect.(type) {
	case TellEffect[W]:
		return r.WithState(r.State().Push(t.output)).Run(ApplyContinuationToEffectResult(t, m.queue, UnitValue))
	}
	return
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

func _[E State[E, S], S any, T any]() {
	// Statically ensure that certain interfaces are implemented correctly
	var _ State[E, S] = StateI[E, S]{}
	var _ TypedEffectTag[S] = GetEffect[S]{}
	var _ TypedEffectTag[Unit] = SetEffect[S]{}
	var _ Interpreter[E, T, StateResult[T, S]] = runState[S, E, T]{}
	var _ InterpreterWithState[runState[S, E, T], S] = runState[S, E, T]{}
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
	return runState[S, E, T]{}.WithState(state).Run(e)
}

type runState[S any, E State[E, S], T any] struct {
	InterpreterWithStateImpl[S, runState[S, E, T], *runState[S, E, T]]
}

func (r runState[_, _, _]) Name() string {
	return "RunState"
}

func (r runState[S, E, T]) Run(e Eff[E, T]) Eff[E, StateResult[T, S]] {
	return RunImpl[E, T, StateResult[T, S]](r, e)
}

func (r runState[S, _, T]) HandlePure(value T) StateResult[T, S] {
	return StateResult[T, S]{
		Value: value,
		State: r.State(),
	}
}

func (r runState[S, E, T]) HandleEffect(effect EffectTag, m Cont[E, T]) (_ Eff[E, StateResult[T, S]]) {
	switch t := m.effect.(type) {
	case GetEffect[S]:
		return r.Run(ApplyContinuationToEffectResult(t, m.queue, r.State()))
	case SetEffect[S]:
		return r.WithState(t.newState).Run(ApplyContinuationToEffectResult(t, m.queue, UnitValue))
	}
	return
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

func _[E Coroutine[E, Y, R], Y any, R any, T any]() {
	// Statically ensure that certain interfaces are implemented correctly
	var _ Coroutine[E, Y, R] = CoroutineI[E, Y, R]{}
	var _ TypedEffectTag[R] = YieldEffect[Y, R]{}
	var _ Interpreter[E, T, CoroutineResult[Y, R, T]] = runCoroutine[E, Y, R, T]{}
}

// DSL implementation for `Coroutine[E, Y, R]`.
type CoroutineI[E Coroutine[E, Y, R], Y any, R any] struct{}

func (_ CoroutineI[E, Y, R]) Yield(output Y) Eff[E, R] {
	return injectEffect[E, R](YieldEffect[Y, R]{output: output})
}

// Effect tag for `Couroutine[E, Y, R].Yield(output Y) Eff[E, R]`.
type YieldEffect[Y any, R any] struct{ output Y }

func (_ YieldEffect[Y, R]) effectResult() R { panic("marker method") }

type CoroutineResume[Y any, R any, T any] func(resumeWith R) CoroutineResult[Y, R, T]

type CoroutineResult[Y any, R any, T any] struct {
	isYield bool
	result  T
	yielded Y
	resume  CoroutineResume[Y, R, T]
}

// Coroutine is done with a result value of type T.
func Done[Y any, R any, T any](value T) CoroutineResult[Y, R, T] {
	return CoroutineResult[Y, R, T]{
		isYield: false,
		result:  value,
	}
}

// Reporting a value of the type Y, and resuming with the value of type R,
// possibly ending with a value of type T.
func Yield[Y any, R any, T any](value Y, resume CoroutineResume[Y, R, T]) CoroutineResult[Y, R, T] {
	return CoroutineResult[Y, R, T]{
		isYield: true,
		yielded: value,
		resume:  resume,
	}
}

// Whether this is the final result of the coroutine.
func (r CoroutineResult[Y, R, T]) IsDone() bool {
	return !r.isYield
}

func (r CoroutineResult[Y, R, T]) Done() (result T, ok bool) {
	return r.result, !r.isYield
}

// Whether this is an intermediate result of the coroutine.
func (r CoroutineResult[Y, R, T]) IsYield() bool {
	return r.isYield
}

func (r CoroutineResult[Y, R, T]) Yielded() (yielded Y, ok bool) {
	return r.yielded, r.isYield
}

// Resume the coroutine with the provided `value`.
func (r CoroutineResult[Y, R, T]) Resume(value R) CoroutineResult[Y, R, T] {
	if !r.isYield {
		panic("cannot resume: coroutine has already completed")
	}
	return r.resume(value)
}

func RunCoroutine[Y any, R any, E Coroutine[E, Y, R], T any](e Eff[E, T]) Eff[E, CoroutineResult[Y, R, T]] {
	return runCoroutine[E, Y, R, T]{}.Run(e)
}

type runCoroutine[E Coroutine[E, Y, R], Y any, R any, T any] struct{}

func (r runCoroutine[E, Y, R, T]) Name() string {
	return "RunState"
}

func (r runCoroutine[E, Y, R, T]) Run(e Eff[E, T]) Eff[E, CoroutineResult[Y, R, T]] {
	return RunImpl[E, T, CoroutineResult[Y, R, T]](r, e)
}

func (r runCoroutine[E, Y, R, T]) HandlePure(value T) CoroutineResult[Y, R, T] {
	return Done[Y, R](value)
}

func (r runCoroutine[E, Y, R, T]) HandleEffect(effect EffectTag, m Cont[E, T]) (_ Eff[E, CoroutineResult[Y, R, T]]) {
	switch t := m.effect.(type) {
	case YieldEffect[Y, R]:
		return Return[E](Yield(t.output, func(resumeWith R) CoroutineResult[Y, R, T] {
			return RunPureOrFail(r.Run(ApplyContinuationToEffectResult(t, m.queue, resumeWith)))
		}))
	}
	return
}
