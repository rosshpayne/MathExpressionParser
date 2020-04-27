package main

import (
	"fmt"
	"strconv"

	"github.com/DynamoGraph/lexer"
	"github.com/DynamoGraph/parser"
	"github.com/DynamoGraph/token"
)

func buildExprGraph(input string) *expression {

	type state struct {
		lvl depthT
		md  bool // true for multiply or division ineffect.
	}

	var (
		tok          *token.Token
		prevTok      *token.Token
		lvl          depthT // precedence level. "(" increments level and adds to graph using extendRight() while ")" decrements level and uses addParent() to existing expression to extend the graph.
		numL         *num   // number node on lhs of operator expression
		numR         *num   // number node on rhs of operator expression
		operandL     bool   // put next INT in numL
		extendRight  bool   // Used when a higher precedence operation detected. Assigns the latest expression to the right operand of the current expression.
		negative     bool   // negative number detected
		cancelRPAREN bool
		multidiv     bool
		opr          operator
		e, en        *expression // "e" points to current expression in graph while "en" is the latest expression to be created and added to the graph using addParent() or extendRight() functions.
		lp           []state
	)

	pushState := func() {
		s := state{lvl: lvl, md: multidiv}
		lp = append(lp, s)
	}

	popState := func() {
		var s state
		s, lp = lp[len(lp)-1], lp[:len(lp)-1]
		lvl = s.lvl
		multidiv = s.md
	}
	// as the parser processes the input left to right it builds a tree (graph) by creating an expression as each operator is parsed and then immediately links
	// it to the previous expression. If the expression is at the same precedence level it links the new expression as the parent of the current expression. In the case
	// of higher precedence operations it links to the right of the current expression (func: extendRight). Walking the tree and evaluating each expression returns the final result.

	fmt.Printf("\n %s \n", input)

	l := lexer.New(input)
	p := parser.New(l)
	operandL = true
	multidiv = true

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

			pushState()
			lvl++

			fmt.Printf("LPAREN.....lvl %d   opr  [%c]  numL= %v\n", lvl, opr, numL)
			// look ahead
			t := p.CurToken
			if t.Type == token.LPAREN {
				continue
			}
			//
			// 3*(3+4      (3+4
			//  ^          ^
			//  opr        opr is nil
			// ^           ^
			// numL        numL is nil
			//
			if opr != 0 {
				// parser is not at the beginning of the stmt e.g (5*3+...
				if numL == nil {

					fmt.Println("IN LPAREN: numL nil make operator only expression and append to current expression")
					//  operands have been consumed into current expression. Take opr and make it parent of the expression.
					//  The next operands and operator to be parsed will be extended to right of this operator only expression.
					en, opr = makeExpr(lvl-1, nil, opr, nil)
					e = e.addParent(en)

				} else {

					// numL has yet to be consumed in expr, this means we have an active * or / operator
					// "(" has interrupted current  operation, so must complete existing extendRight
					fmt.Println("IN LPAREN: numL not nil make numL only expression ")
					en, opr = makeExpr(lvl-1, numL, opr, nil)

					if e == nil {
						fmt.Println("IN LPAREN: set expr as first time ")
						e, en = en, nil

					} else {
						fmt.Println("IN LPAREN: EXTEND RIGHT...on current expr ")
						e, lvl = e.extendRight(en, lvl)
					}
				}
				operandL = true
				extendRight = true
				numL, numR = nil, nil

			}

		case token.RPAREN:

			popState()

			fmt.Println("RPAREN  lvl ", lvl)
			// peek and next token
			t := p.CurToken
			if (t.Type == token.MULTIPLY || t.Type == token.DIVIDE) && !multidiv {
				//	delay impact of RPAREN until next RPAREN. Any "+", "-" will be added to this level which is OK
				//  as these operations do not impact the precedence. Any "*","/" will
				//  extend-right as usual and increase the level as a result. (see extendRight)
				lvl++
				cancelRPAREN = true
			}
			fmt.Println("RPAREN  lvl ", lvl)

		case token.INT:

			i, err := strconv.Atoi(tok.Literal)
			if err != nil {
				//l.Err = fmt.Errorf("Failed to convert to int, %w", err)
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
			fmt.Printf("********* in INT %d : opr = [%c]. lvl = %d\n", i, opr, lvl)
			if opr == '+' || opr == '-' {
				//
				tok := p.CurToken // which is really a peekToken because we do a NextToken after setting tok in top of for loop
				if tok.Type == token.MULTIPLY || tok.Type == token.DIVIDE {
					//
					// High precedence operaton - create node (expression) and attach to graph in preparation for future extendRight node(s).
					//
					if extendRight {
						en, opr = makeExpr(lvl, numL, opr, nil)
						if e == nil {
							e, en = en, nil
						} else {
							e, lvl = e.extendRight(en, lvl)
							extendRight = false
						}

					} else if numL == nil {
						// add operator only node to graph - no left, right operands. addParent will attach left, and future ExtendRIght will attach right.
						en, opr = makeExpr(lvl, nil, opr, nil)
						e = e.addParent(en)

					} else {
						// make expr for existing numL and opr
						en, opr = makeExpr(lvl, numL, opr, nil)
						if e == nil {
							e, en = en, nil
						} else {
							e = e.addParent(en)
						}
					}
					fmt.Println("HIGHER PRECEDENCE ....  ", i, lvl)
					// all higher precedence operations or explicit (), perform an "extendRight" to create a new branch in the graph.
					extendRight = true
					// new branches begin with a left operand
					operandL = true
				}
			}

			if operandL {

				numL = &num{i: i}
				operandL = false
				fmt.Println("Left NUM", i, "  lvl: ", lvl)

			} else {

				numR = &num{i: i}
				fmt.Println("Right NUM ", i, "  lvl: ", lvl)
				en, opr = makeExpr(lvl, numL, opr, numR)
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
						fmt.Printf("**** extendRight on  %c lvl %d     child: %c  %s lvl. %d \n\n", e.opr, e.depth, en.opr, en.name, en.depth)
						// higher precedence operator or a ( has occured - create new branch to right of current expression
						e, lvl = e.extendRight(en, lvl)
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
			if cancelRPAREN {
				lvl--
				cancelRPAREN = false
			}
			multidiv = false

		case token.MINUS:
			// is it a negative sign or a minus sign?
			ptok := prevTok.Type
			if ptok == token.LPAREN || ptok == token.MULTIPLY || ptok == token.PLUS || ptok == token.MINUS || ptok == token.DIVIDE {
				negative = true
			} else {
				opr = MINUS
				if cancelRPAREN {
					lvl--
					cancelRPAREN = false
				}
			}
			multidiv = false

		case token.MULTIPLY:

			opr = MULTIPLY
			multidiv = true

		case token.DIVIDE:

			opr = DIVIDE
			multidiv = true
		}
		if tok.Type == token.EOF {
			break
		}

	}
	return findRoot(e)
}
