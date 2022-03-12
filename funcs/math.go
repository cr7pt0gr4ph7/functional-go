package funcs

import (
	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Complex | constraints.Integer | constraints.Float
}

func Add[T Numeric | ~string](x T, y T) T {
	return x + y
}

func Sub[T Numeric](x T, y T) T {
	return x - y
}

func Mul[T Numeric](x T, y T) T {
	return x * y
}

func Div[T Numeric](x T, y T) T {
	return x / y
}

func Plus[A Numeric | ~string](summand A) func(in A) A {
	return func(in A) A {
		return in + summand
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
