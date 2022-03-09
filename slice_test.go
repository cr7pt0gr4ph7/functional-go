package functional

func ExampleSlice_FoldLeft() {
	s := Slice[int]{1, 2, 3, 4}
	r := FoldLeft(s, 10, func(elem int, state int) int {
		fmt.Println(elem, state)
		return elem + state
	})
	fmt.Println("=", r)
	// Output:
	// 1 10
	// 2 11
	// 3 13
	// 4 16
	// = 20
}
