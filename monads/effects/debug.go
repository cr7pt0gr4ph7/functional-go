package effects

// Debugging utilities
//
// This file contains debugging helpers that aren't used during normal execution.
// Almost no code outside of this file should depend on the methods in this file,
// to enable this file to be omitted from the build as an optimization.
//
// Methods that ARE referenced by outside code should have a stub provided in
// `debug_stubs.go`.

import (
	"fmt"
	"reflect"
)

//
// Debug helpers for effects interpreters
//

type logType struct{}

var log logType

func (_ logType) OnRunEffect(runnerName string, args ...any) {
	fmt.Println(append([]any{runnerName}, args...)...)
}

//
// Debug helpers for Eff[E, A]
//

func (p Pure[E, A]) String() string {
	return fmt.Sprintf("Pure(%v)", p.value)
}

func (c Cont[E, B]) String() string {
	return fmt.Sprintf("Cont(%v -> %v)", effectTagToString(c.effect), c.queue)
}

func effectTagToString(tag EffectTag) string {
	if t, ok := tag.(debuggableEffectTag); ok {
		return t.effectTagString()
	}
	return fmt.Sprintf("%v(%v)", reflect.TypeOf(tag), tag)
}

type debuggableEffectTag interface {
	effectTagString() string
}

//
// Debug helpers for evalQueue[E, A, B]
//

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

func (l identQ[E, A]) getType(typ typeToRetrieve, r recursionEnum) reflect.Type {
	switch typ {
	case effectsType:
		return typeOf[E]()
	case leftType:
		return typeOf[A]()
	case middleType:
		return nil
	case rightType:
		return typeOf[A]()
	default:
		return nil
	}
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

func (t runQ[E, B, C]) getType(typ typeToRetrieve, r recursionEnum) reflect.Type {
	switch typ {
	case effectsType:
		return typeOf[E]()
	case leftType:
		if r == recurseAlways || r == recurseIfNeeded {
			return getType(t.wrapped, leftType, r)
		}
		return nil
	case middleType:
		if r == recurseAlways {
			return getType(t.wrapped, rightType, r)
		}
		return typeOf[B]()
	case rightType:
		return typeOf[C]()
	default:
		return nil
	}
}

func (t nodeQErased[E, B]) getType(typ typeToRetrieve, r recursionEnum) reflect.Type {
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

func (l identQ[E, A]) String() string {
	return fmt.Sprintf("{PassThru: %v}", l.getType(leftType, recurseIfNeeded))
}

func (l leafQ[E, A, B]) String() string {
	if len(l.debugTag) > 0 {
		return fmt.Sprintf("{%v: %v => %v}", l.debugTag, l.getType(leftType, recurseIfNeeded), l.getType(rightType, recurseIfNeeded))
	} else {
		return fmt.Sprintf("{%v => %v}", l.getType(leftType, recurseIfNeeded), l.getType(rightType, recurseIfNeeded))
	}
}

func (t runQ[E, B, C]) String() string {
	return fmt.Sprintf("[%v ==> %v]", t.wrapped, t.debugTag)
}

func (t nodeQErased[E, B]) String() string {
	return fmt.Sprintf("(%v => %v)", t.left, t.right)
}
