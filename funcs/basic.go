package funcs

func Nop() {
}

func Ignore[T any](_ T) {
}

func Identity[T any](arg T) T {
	return arg
}

func Panic(reason any) func() {
	return func() {
		panic(reason)
	}
}

func Compose[A any, B any, C any](f func(arg A) B, g func(arg B) C) func(arg A) C {
	return func(arg A) C {
		return g(f(arg))
	}
}

func ComposeErr[A any, B any, C any](f func(arg A) (B, error), g func(arg B) (C, error)) func(arg A) (C, error) {
	return func(arg A) (_ C, err error) {
		if b, e := f(arg); e != nil {
			err = e
			return
		} else {
			return g(b)
		}
	}
}
