package eval

func ExampleEval() {
	x0 := Now(42)
	x1 := Map(x0, func(x int) int { return x / 3 })
	x2 := FlatMap(x1, func(x int) Eval[float32] {
		return Always(func() float32 {
			return float32(x) * 0.7
		})
	})
	fmt.Println(x2.Value())
	// Output: 9.8
}
