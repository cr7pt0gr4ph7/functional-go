package effects

type Interpreter[E any, A any, B any] interface {
	Name() string
	Run(e Eff[E, A]) Eff[E, B]
	HandlePure(value A) B
	HandleEffect(effect EffectTag, m Cont[E, A]) Eff[E, B]
}

type InterpreterWithState[Self any, S any] interface {
	State() S
	WithState(newState S) Self
}

type stateForDebug interface {
	stateForDebug() any
}

type InterpreterWithStateImpl[S any, Self any, SelfPtr interface {
	*Self
	unsafeSetState(newState S)
}] struct{ state S }

func (r InterpreterWithStateImpl[S, Self, SelfPtr]) stateForDebug() any {
	return r.state
}

func (r *InterpreterWithStateImpl[S, Self, SelfPtr]) unsafeSetState(newState S) {
	r.state = newState
}

func (r InterpreterWithStateImpl[S, Self, SelfPtr]) State() S {
	return r.state
}

func (r InterpreterWithStateImpl[S, Self, SelfPtr]) WithState(newState S) Self {
	var newSelf Self
	SelfPtr(&newSelf).unsafeSetState(newState)
	return newSelf
}

func RunImpl[E any, A any, B any](r Interpreter[E, A, B], e Eff[E, A]) Eff[E, B] {
	if d, ok := r.(stateForDebug); ok {
		log.OnRunEffect(r.Name(), d.stateForDebug(), e)
	} else {
		log.OnRunEffect(r.Name(), e)
	}

	switch m := e.EffImpl.(type) {
	case Pure[E, A]:
		return newPure[E](r.HandlePure(m.value))
	case Cont[E, A]:
		result := r.HandleEffect(m.effect, m)
		if result.EffImpl == nil {
			// Hack: Use the default value of Eff[_, _] as a sentinel
			//       to signal that the effect has not been handled yet.
			return ForwardEffect(m, r.Run, r.Name())
		}
		return result
	default:
		panic("unreachable")
	}
}
