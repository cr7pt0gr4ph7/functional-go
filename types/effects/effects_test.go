package effects

type MyEffects interface {
	Reader[MyEffects, float32]
	Writer[MyEffects, string]
	State[MyEffects, int]
}

type MyEffectsI struct {
	ReaderI[MyEffects, float32]
	WriterI[MyEffects, string]
	StateI[MyEffects, int]
}

func ExampleMyEffects() {
	var e MyEffectsI

	x0 := Return[MyEffects](42)
	x1 := Map(x0, func(x int) int { return x / 2 })

	fmt.Println(x0, x1)

	y0 := Chain(
		e.Tell("Hello, "),
		e.Tell("environment"),
		e.Tell("!\n"),
	)

	yr := RunPureOrFail(RunWriter[string](y0))

	fmt.Println(yr)
}
