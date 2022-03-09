package funcs

import (
	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Complex | constraints.Integer | constraints.Float
}

func Plus[A Numeric | ~string](add A) func(in A) A {
	return func(in A) A {
		return in + add
	}
}

func Times[A Numeric](factor A) func(in A) A {
	return func(in A) A {
		return in * factor
	}
}

func DividedBy[A Numeric](quotient A) func(in A) A {
	return func(in A) A {
		return in / quotient
	}
}
