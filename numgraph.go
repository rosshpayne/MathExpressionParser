package main

import (
	"fmt"
	"strconv"

	"github.com/DynamoGraph/lexer"
	"github.com/DynamoGraph/parser"
	"github.com/DynamoGraph/token"
)

const (
	LPAREN uint8 = 1
)

func buildExprGraph(input string) *expression {

	type state struct {
		opr operator
	}

	var (
		tok         *token.Token
		prevTok     *token.Token
		numL        *num // number node on lhs of operator expression
		numR        *num // number node on rhs of operator expression
		operandL    bool // put next INT in numL
		extendRight bool // Used when a higher precedence operation detected. Assigns the latest expression to the right operand of the current expression.
		negative    bool // negative number detected
		opr         operator
		opr_        operator    // saved operator from state
		e, en       *expression // "e" points to current expression in graph while "en" is the latest expression to be created and added to the graph using addParent() or extendRight() functions.
		lp          []state
	)

	pushState := func() {
		s := state{opr: opr}
		lp = append(lp, s)
		fmt.Printf("\n================================================== PUSH: len,opr_ %d %c \n", len(lp), s.opr)
	}

	popState := func() {
		var s state
		s, lp = lp[len(lp)-1], lp[:len(lp)-1]
		opr_ = s.opr
		fmt.Printf("\n================================================== POP: len,opr_  %d %c \n", len(lp), s.opr)

	}

	// as the parser processes the input left to right it builds a tree (graph) by creating an expression as each operator is parsed and then immediately links
	// it to the previous expression. If the expression is at the same precedence level it links the new expression as the parent of the current expression. In the case
	// of higher precedence operations it links to the right of the current expression (func: extendRight). Walking the tree and evaluating each expression returns the final result.

	fmt.Printf("\n %s \n", input)

	l := lexer.New(input)
	p := parser.New(l)
	operandL = true

	// TODO - initial full parse to validate left and right parenthesis match

	for {

		prevTok = tok
		tok = p.CurToken
		p.NextToken()
		fmt.Printf("\ntoken: %s\n", tok.Type)

		switch tok.Type {
		case token.EOF:
			break
		case token.LPAREN:
			//
			// LPAREN is represented in the graph by a "NULL" expression (node) consisting of operator "+" and left operand of 1.
			//
			// 3*(3+4      (3+4
			//  ^          ^
			//  opr        opr is nil
			// ^           ^
			// numL        numL is nil
			//
			// save the current state as we will need it when the associated closing param is evaluated
			//
			pushState()
			//
			// add any expression that has not already been added to the graph
			//
			if opr != 0 {
				// parser is not at the beginning of the stmt e.g (5*3+...
				if numL == nil {

					fmt.Println("IN LPAREN: numL nil make operator only expression and append to current expression")
					//  operands have been consumed into current expression. Take opr and make it parent of the expression.
					//  The next operands and operator to be parsed will be extended to right of this operator only expression.
					en, opr = makeExpr(nil, opr, nil)
					e = e.addParent(en)

				} else {

					// numL has yet to be consumed in expr, this means we have an active * or / operator
					// "(" has interrupted current  operation, so must complete existing extendRight
					fmt.Println("IN LPAREN: numL not nil make numL only expression ")
					en, opr = makeExpr(numL, opr, nil)

					if e == nil {
						e, en = en, nil

					} else {
						e = e.extendRight(en)
					}
				}
			}
			//
			// add a NULL expression representing the "("
			//
			en = &expression{left: &num{i: 0}, opr: '+', right: nil, name: "LPAREN", id: LPAREN}
			if e == nil {
				e, en = en, nil

			} else {
				e = e.extendRight(en)
			}
			operandL = true
			extendRight = true
			numL, numR = nil, nil
			fmt.Println("END LPAREN: e.name : ", e.name)

		case token.RPAREN:

			fmt.Println("RPAREN  ")

			popState()

			// navigate current expression e, up to next LPARAM expression
			for e = e.parent; e.parent != nil && e.id != LPAREN; e = e.parent {
			}
			//
			if e.parent != nil && e.parent.id != LPAREN {
				// opr_ represents the operator that existed at the associated "(". Sourced from state.
				if opr_ == '/' || opr_ == '*' || opr_ == 0 {
					fmt.Println("opr_ is adjusting to e.parent ", opr_)
					e = e.parent
				}
			}
			if e.parent != nil {
				fmt.Println("+finalis e.name, parent : ", e.name, e.parent.name)
			} else {
				fmt.Println("final e.name , parent: ", e.name, "nil")
			}

		case token.INT:

			i, err := strconv.Atoi(tok.Literal)
			if err != nil {
				break
			}
			if negative {
				i *= -1
				negative = false
			}
			//
			// before assigning numL or numR, check for higher precedence operation
			//
			// 2+5*7		3+4*7       )*5+4*7
			//   ^            ^             ^
			// ^            ^             ^
			// numL        numL          numR (numL nil)
			//
			fmt.Printf("********* in INT %d : opr = [%c]\n", i, opr)
			if opr == '+' || opr == '-' {
				//
				tok := p.CurToken // which is really a peekToken because we do a NextToken after setting tok in top of for loop
				if tok.Type == token.MULTIPLY || tok.Type == token.DIVIDE {
					//	hp = true
					//
					// High precedence operaton - create node (expression) and attach to graph in preparation for future extendRight node(s).
					//
					if extendRight {
						en, opr = makeExpr(numL, opr, nil)
						if e == nil {
							e, en = en, nil
						} else {
							e = e.extendRight(en)
							extendRight = false
						}

					} else if numL == nil {
						// add operator only node to graph - no left, right operands. addParent will attach left, and future ExtendRIght will attach right.
						en, opr = makeExpr(nil, opr, nil)
						e = e.addParent(en)

					} else {
						// make expr for existing numL and opr
						en, opr = makeExpr(numL, opr, nil)
						if e == nil {
							e, en = en, nil
						} else {
							e = e.addParent(en)
						}
					}
					// all higher precedence operations or explicit (), perform an "extendRight" to create a new branch in the graph.
					extendRight = true
					// new branches begin with a left operand
					operandL = true
				}
			}

			if operandL {

				numL = &num{i: i}
				operandL = false
				fmt.Println("Left NUM ", i)

			} else {

				numR = &num{i: i}
				fmt.Println("Right NUM ", i)
				en, opr = makeExpr(numL, opr, numR)
				if e == nil {
					e, en = en, nil
				}
				// consumed following values, so reset them
				numL, numR = nil, nil
				// addParent is the default operation to extend the graph, which requires a numR only
				operandL = false
				// do not extend or add expression until we have an "e" and "en" expression.
				if en != nil {
					if extendRight {
						fmt.Printf("**** extendRight on  %c    child: %c  %s \n\n", e.opr, en.opr, en.name)
						e = e.extendRight(en)
						// addParent is the default method to extend the graph, so make extendRight false. Must be explicitly set to true
						// when the correct scenario occurs i.e. immediately after a ( or higher precedence operation detected
						extendRight = false
					} else {
						e = e.addParent(en)
					}
				}
			}

		case token.PLUS:

			opr = PLUS

		case token.MINUS:
			// is it a negative sign or a minus sign?
			ptok := prevTok.Type
			if ptok == token.LPAREN || ptok == token.MULTIPLY || ptok == token.PLUS || ptok == token.MINUS || ptok == token.DIVIDE {
				negative = true
			} else {
				opr = MINUS
			}

		case token.MULTIPLY:

			opr = MULTIPLY

		case token.DIVIDE:

			opr = DIVIDE

		}
		if tok.Type == token.EOF {
			break
		}

	}
	return findRoot(e)
}
