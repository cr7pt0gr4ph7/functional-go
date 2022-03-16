package effects

import (
	"fmt"
	"github.com/cr7pt0gr4ph7/functional-go/types/list"
)

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

func ExampleMyEffects2() {
	fmt.Println("Hello, playground")
	var e MyEffectsI

	x0 := Return[MyEffects](42)
	x1 := Map(x0, func(x int) int { return x / 2 })

	fmt.Println(x1)

	y0 := Chain(
		e.Tell("Hello, ").And(e.Tell("external ")),
		// e.Yield(10).Discard(),
		e.Get().Discard(),
		e.Set(10),
		e.Get().And(e.Tell("effects ")).Discard(),
		e.Tell("environment"),
		// e.Yield(20).Discard(),
		e.Get().Discard(),
		e.Set(20),
		e.Get().Discard(),
		e.Tell("!\n"),
	)

	fmt.Println(y0)

	yr := RunPureOrFail(RunState(5, RunWriter[string](y0)))

	fmt.Println("State:", yr.State)
	fmt.Println("Written:", yr.Value.Written)
	fmt.Println("Result:", yr.Value.Value)

	yr1 := RunPureOrFail(RunWriter[string](RunState(5, y0)))

	fmt.Println("State:", yr1.Value.State)
	fmt.Println("Written:", yr1.Written)
	fmt.Println("Result:", yr1.Value.Value)

	yr2 := RunPureOrFail(RunWriterReverse[string](list.Empty[string](), RunState(5, y0)))

	fmt.Println("State:", yr2.Value.State)
	fmt.Println("Written:", yr2.Written)
	fmt.Println("Result:", yr2.Value.Value)
}
