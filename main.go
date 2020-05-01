package main

import (
	"fmt"
	"strconv"

	"github.com/satori/go.uuid"
)

type literalT uint8

const (
	num_ literalT = iota
	str
	bl
)

type funcT uint8

const (
	eq funcT = iota
	le
	lt
	gt
	ge
)

type uidT = uuid.UUID

type uidS []uidT

type litValI interface {
	lit() uidS
}

type predNameT string

type valS string

func (s *valS) lit() {}

type valN float64

func (n *valN) lit() {}

type valBl bool

func (b *valBl) lit() {}

type Function func(p predNameT, litVal string) uidS

type operator byte

const (
	PLUS     operator = '+'
	MINUS    operator = '-'
	MULTIPLY operator = '*'
	DIVIDE   operator = '/'
)

func walk(e operand) {

	if e, ok := e.(*expression); ok {

		walk(e.left)
		walk(e.right)

		switch e.opr {
		case PLUS:

			e.result = e.left.getResult() + e.right.getResult()

		case MINUS:

			e.result = e.left.getResult() - e.right.getResult()

		case MULTIPLY:

			e.result = e.left.getResult() * e.right.getResult()

		case DIVIDE:

			e.result = e.left.getResult() / e.right.getResult()

		}
		fmt.Printf("Result: %c %g\n", e.opr, e.result)
	}

}

func findRoot(e *expression) *expression {

	for e.parent != nil {
		e = e.parent
	}
	return e
}

// operand interface.
// So far type num (integer), expression satisfy, but this can of course be extended to floats, complex numbers, functions etc.
type operand interface {
	getParent() *expression
	type_() string
	printName() string
	getResult() float64
}

type expression struct { // expr1 and expr2     expr1 or expr2       exp1 or (expr2 and expr3). (expr1 or expr2) and expr3
	id     uint8    // type of expression. So far used only to identify the NULL expression, representing the "(" i.e the left parameter or LPARAM in a mathematical expression
	name   string   // optionally give each expression a name. Maybe useful for debugging purposes.
	result float64  // store result of "left operator right. Walking the graph will interrogate each operand for its result.
	left   operand  //
	opr    operator // for Boolean: AND OR NOT  For mathematical: +-/*
	right  operand  //
	parent *expression
}

func (e *expression) getParent() *expression {
	return e.parent
}
func (e *expression) type_() string {
	return "expression"
}
func (e *expression) printName() string {
	return e.name
}

func (e *expression) getResult() float64 {
	return e.result
}

// integer
type num struct {
	parent *expression
	i      int
}

func (n *num) getParent() *expression {
	return n.parent
}

func (n *num) type_() string {
	return "num"
}

func (n *num) getResult() float64 {
	return float64(n.i)
}

func (n *num) printName() string {
	if n == nil {
		return "nil"
	} else {
		return fmt.Sprintf("%d", n.i)
	}
}

func makeExpr(l operand, op operator, r operand) (*expression, operator) {

	e := &expression{left: l, opr: op, right: r}

	// remember: nil interfaces means the type component is nil not necessarily the value component.
	// if a nil numL is passed to makeExpr, the type component is set (operand) but the value (concrete type) is nil.
	// so to check the interface is nil you must check the value is also nil, as below.
	if x, ok := e.left.(*num); ok {
		if x != nil {
			x.parent = e
		}
	}
	if x, ok := e.right.(*num); ok {
		if x != nil {
			x.parent = e
		}
	}
	if l != nil && r != nil {
		if l_, ok := l.(*num); ok {
			if r_, ok := r.(*num); ok {
				ln := "nil "
				if l_ != nil {
					ln = strconv.Itoa(l_.i) + " "
				}
				rn := " nil"
				if r_ != nil {
					rn = " " + strconv.Itoa(r_.i)
				}
				e.name = "[ " + ln + string(op) + rn + "]"
			}
		}
	} else if l == nil && r == nil {
		e.name = "[" + "nil " + string(op) + " nil" + "]"
	} else {
		if l != nil {
			if x, ok := l.(*num); ok {
				v := "nil"
				if x != nil {
					v = strconv.Itoa(x.i) + " "
				}
				e.name = "[ " + v + string(op) + " nil" + "]"
			}
		}
		if r != nil {
			if x, ok := r.(*num); ok {
				v := "nil"
				if x != nil {
					v = " " + strconv.Itoa(x.i)
				}
				e.name = "[ " + "nil " + string(op) + v + "]"
			}
		}
	}
	fmt.Printf("\n****   MakeExpr %s\n", e.name)

	return e, 0
}

// ExtendRight for Higher Precedence operators or open braces - parsed:   *,/, (
// c - current op node, n is the higer order op we want to extend right
func (c *expression) extendRight(n *expression) *expression {

	c.right = n
	n.parent = c

	fmt.Printf("++++++++++++++++++++++++++ extendRight  -  FROM %s  -> %s  \n", c.name, n.name)
	return n
}

func (c *expression) addParent(n *expression) *expression {
	//
	if c.parent != nil {
		//  current node must now point to the new node being added, and similar the new node must point back to the current node.
		c.parent.right = n
		n.parent = c.parent
	}
	// set old parent to new node
	c.parent = n
	n.left = c

	fmt.Printf("\n++++++++++++++++++++++++++ addParent  %s on %s \n\n", n.name, c.name)
	return n
}

type function struct {
	parent *expression
	name   string
	f      func(string, litValI) bool // eq
	arg1   string                     // facet "close", "relative"
	//arg2Typ literalT                    // boolean
	arg2Val litValI // true
}

func (f *function) oper() {}

func (f *function) getParent() *expression {
	return f.parent
}

func (f *function) type_() string {
	return "func"
}

func (f *function) getResult() bool {
	return f.f(f.arg1, f.arg2Val)
}

func (f *function) printName() string {
	if f == nil {
		return "nil"
	}
	return f.name
}

func main() {}
