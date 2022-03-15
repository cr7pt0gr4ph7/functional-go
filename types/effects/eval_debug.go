package effects

import (
	"fmt"
	"reflect"
)

type typeToRetrieve byte

const (
	effectsType typeToRetrieve = iota
	leftType
	middleType
	rightType
)

type recursionEnum byte

const (
	recurseIfNeeded recursionEnum = iota
	recurseAlways
	recurseNever
)

type debuggableTreeNode interface {
	getType(t typeToRetrieve, r recursionEnum) reflect.Type
	String() string
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf(new(T)).Elem()
}

func getType(node evalTreeNode, typ typeToRetrieve, r recursionEnum) reflect.Type {
	if node == nil {
		return nil
	}
	return node.(debuggableTreeNode).getType(typ, r)
}

func (l leafQ[E, A, B]) getType(typ typeToRetrieve, r recursionEnum) reflect.Type {
	switch typ {
	case effectsType:
		return typeOf[E]()
	case leftType:
		return typeOf[A]()
	case middleType:
		return nil
	case rightType:
		return typeOf[B]()
	default:
		return nil
	}
}

func (t nodeQ[E, A, X, B]) getType(typ typeToRetrieve, r recursionEnum) reflect.Type {
	switch typ {
	case effectsType:
		return typeOf[E]()
	case leftType:
		if r == recurseAlways {
			return getType(t.left, leftType, r)
		}
		return typeOf[A]()
	case middleType:
		if r == recurseAlways {
			return getType(t.left, rightType, r)
		}
		return typeOf[X]()
	case rightType:
		if r == recurseAlways {
			return getType(t.right, rightType, r)
		}
		return typeOf[B]()
	default:
		return nil
	}
}

func (t nodeQ2[E, B]) getType(typ typeToRetrieve, r recursionEnum) reflect.Type {
	switch typ {
	case effectsType:
		return typeOf[E]()
	case leftType:
		if r == recurseAlways || r == recurseIfNeeded {
			return getType(t.left, leftType, r)
		}
		return nil
	case middleType:
		if r == recurseAlways || r == recurseIfNeeded {
			return getType(t.left, rightType, r)
		}
		return nil
	case rightType:
		if r == recurseAlways {
			return getType(t.right, rightType, r)
		}
		return typeOf[B]()
	default:
		return nil
	}
}

func (l leafQ[E, A, B]) String() string {
	return fmt.Sprintf("{%v => %v}", l.getType(leftType, recurseIfNeeded), l.getType(rightType, recurseIfNeeded))
}

func (t nodeQ[E, A, X, B]) String() string {
	return fmt.Sprintf("(%v => %v)", t.left, t.right)
}

func (t nodeQ2[E, B]) String() string {
	return fmt.Sprintf("[%v => %v]", t.left, t.right)
}
