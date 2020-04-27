package main

import (
	"fmt"
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
	eq funcT = 1 << iota
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

// {
//   data(func: eq(name, "Alice")) {
//     friend @facets(eq(close, true) AND eq(relative, true)) @facets(relative) { # filter close friends in my relation
//       name
//     }
//   }
// }
// for each uid
//.   execute filterExpr(uid)
//		walk query tree
//  		when it gets to func. fetch from dynamo, given uid and what facets are required
//          run function with values from db. Get true/false. If true save uid in uids. If false ignore and go on to next.
// for each uid in uids
//       walk the parse tree
//         and output data held in the tree from dynamo
//

type operator byte

const (
	PLUS     operator = '+'
	MINUS    operator = '-'
	MULTIPLY operator = '*'
	DIVIDE   operator = '/'
)

type depthT = int8

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

// operand interface. type num, expression satisfy.
type operand interface {
	getParent() *expression
	type_() string
	printName() string
	getResult() float64
}

type expression struct { // expr1 and expr2     expr1 or expr2       exp1 or (expr2 and expr3). (expr1 or expr2) and expr3
	depth  depthT   // precedence level. TODO: rename to something more suitable.
	name   string   // optionally give each expression a name. Maybe useful for debugging purposes.
	result float64  // store result of "left.int operator right.int". Walking the graph will interrogate each operand (interface) for its result.
	left   operand  // eq(close, true) , order  value comes from  uid in the set of uids belonging to the predicate.
	opr    operator // for Boolean: AND OR NOT  For mathematical: +-/*
	right  operand  // eq(relative, true)
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

func makeExpr(d depthT, l operand, op operator, r operand) (*expression, operator) {

	e := &expression{depth: d, left: l, opr: op, right: r}
	fmt.Printf("MakeExpr depth  %d opr %c  %v %v\n", e.depth, op, e.left, e.right)

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
	return e, 0
}

// ExtendRight for Higher Precedence operators or open braces - parsed:   *,/, (
// c - current op node, n is the higer order op we want to extend right
func (c *expression) extendRight(n *expression, lvl depthT) (*expression, depthT) {

	c.right = n
	//n.depth = lvl + 1
	n.depth = c.depth + 1
	n.parent = c
	fmt.Printf("ExtendRight......%c-%d  %c-%d\n", c.opr, c.depth, n.opr, n.depth)
	return n, n.depth
}

// addParent - add expression (argument), to expression (method receiver) as a parent, if it is at a suitable level in the precedence hierarchy.
// Otherwise recursively walk the graph upwards until we get to an expression with the correct precedence level.
func (c *expression) addParent(n *expression) *expression {
	//
	// based on depth (precedence level) of expression n and c, walk up the tree to find a suitable expression to append n to.
	//
	ediff := c.depth - n.depth
	if ediff > 0 {
		// move to next paranthesis level, noting that the expression may have no parent.
		if c.parent != nil {
			fmt.Println("addParent ===== lvl ", c.depth, n.depth)
			return c.parent.addParent(n)
		}
	}
	//
	// At the correct paranthesis level. Now add the new expression, n, as parent.
	//
	if c.parent != nil {
		// as with ExtendRight(), the parent of the current node, if it exists at this level,
		//  must now point to the new node being added, and similar the new node must point back to the current node.
		c.parent.right = n
		n.parent = c.parent
	}
	// set old parent to new node
	c.parent = n
	n.left = c
	fmt.Printf("addParent.....new parent %c-%d on %c-%d\n", n.opr, n.depth, c.opr, c.depth)

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
