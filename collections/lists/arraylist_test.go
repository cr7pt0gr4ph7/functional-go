package lists

import (
	"fmt"
)

func ExampleArrayList() {
	var l ArrayList[int]
	fmt.Println(l)

	l.Add(10)
	l.Add(20)
	l.Add(30)
	fmt.Println(l)

	l.Clear()
	fmt.Println(l)
}
