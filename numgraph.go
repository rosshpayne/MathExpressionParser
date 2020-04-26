package main

import (
	"fmt"
	"strconv"

	"github.com/DynamoGraph/lexer"
	"github.com/DynamoGraph/parser"
	"github.com/DynamoGraph/token"
)

func numGraph(input string) *expression {

	var (
		tok         *token.Token
		prevTok     *token.Token
		lvl         depthT // precedence level. "(" increments level and adds to graph using extendRight() while ")" decrements level and uses addParent() to existing expression to extend the graph.
		numL        *num   // number node on lhs of operator expression
		numR        *num   // number node on rhs of operator expression
		operandL    bool   // put next INT in numL
		extendRight bool   // Used when a higher precedence operation detected. Assigns the latest expression to the right operand of the current expression.
		negative    bool   // negative number detected
		ihpSet      bool
		opr         operator
		e, en       *expression // "e" points to current expression in graph while "en" is the latest expression to be created and added to the graph using AddParent() or ExtendRight() functions.
		lp          []depthT
	)
	// as the parser processes left to right it builds a tree (graph) by creating an expression as each operator is parsed and then immediately links
	// it to the previous expression when its not the first operator to be parsed. For equal or lower precedence it links the new expression as the parent of the previous one or
	// for higher precedence operations to the right of the previous expression (extendRight). By walking the tree lower left to right the expressions are evaluated
	// and summed together to get a final result.
	fmt.Println()
	fmt.Println(input)
	fmt.Println()
	l := lexer.New(input)
	p := parser.New(l)
	operandL = true

	// TODO - initial single parse to validate parenthesis match

	for {

		prevTok = tok
		tok = p.CurToken

		p.NextToken()

		fmt.Println("\ntoken: ", tok.Type)
		fmt.Println()
		switch tok.Type {
		case token.EOF:
			break
		case token.LPAREN:

			// push lvl onto lp
			lp = append(lp, lvl)
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
				// "(" has interrupted numL
				if numL == nil {

					fmt.Println("IN LPAREN: numL nil make operator only expression and append to current expression")
					//  operands have been consumed into current expression. Take opr and make it parent of the expression.
					//  The next operands and operator to be parsed will be extended to right of this operator only expression.
					en = makeExpr(lvl-1, nil, opr, nil)
					e = e.AddParent(en)

				} else {

					// numL has yet to be consumed in expr, this means we have an active * or / operator ie. ihpSet true
					// "(" has interrupted current ihpSet operation, so must complete existing extendRight
					fmt.Println("IN LPAREN: numL not nil make numL only expression ")
					en = makeExpr(lvl-1, numL, opr, nil)

					if e == nil {
						fmt.Println("IN LPAREN: set expr as first time ")
						e, en = en, nil

					} else {
						fmt.Println("IN LPAREN: EXTEND RIGHT...on current expr ")
						e, lvl = e.ExtendRight(en, lvl)
					}
				}
				ihpSet = false
				operandL = true
				extendRight = true
				numL, numR = nil, nil
				opr = 0

			}

		case token.RPAREN: // )

			// pop lvl from lp
			lvl, lp = lp[len(lp)-1], lp[:len(lp)-1]
			fmt.Println("in RPAREN: lvl  ", lvl)
			if ihpSet {
				lvl--
				ihpSet = false
				fmt.Println("in RPAREN:  ihpSet ", lvl)
			}
			t := p.CurToken
			if t.Type == token.MULTIPLY || t.Type == token.DIVIDE {
				ihpSet = true
				lvl++
				fmt.Println("in RPAREN:  set ihpSet ", lvl)
			}

			fmt.Println("in RPAREN: lvl--", lvl)

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
				// check if next opr has higher precedence. If so unary expression with previous numL (not current INT) and append to existing lvl.
				//  Next operation with the higher precedence will extendRight from current node.
				//  Set OperandL = true to start the creation an expression with numL * numR
				//
				tok := p.CurToken // which is really a peekToken because we do a NextToken after setting tok in top of for loop
				if tok.Type == token.MULTIPLY || tok.Type == token.DIVIDE {
					//
					// create node (expression) and attach to graph in preparation for future extendRight node(s).
					//
					if extendRight {
						en = makeExpr(lvl, numL, opr, nil)
						if e == nil {
							e, en = en, nil
						} else {
							e, lvl = e.ExtendRight(en, lvl)
							extendRight = false
						}

					} else if numL == nil {
						// add operator only node to graph - no left, right operands. AddParent will attach left, and future ExtendRIght will attach right.
						en = makeExpr(lvl, nil, opr, nil)
						e = e.AddParent(en)

					} else {
						// make expr for existing numL and opr
						en = makeExpr(lvl, numL, opr, nil)
						if e == nil {
							e, en = en, nil
						} else {
							e = e.AddParent(en)
						}

					}
					fmt.Println("HIGHER PRECEDENCE ....  ", i, lvl)
					// all higher precedence operations, ihp or explicit (), perform an "extendRight" to create a separate path in the graph.
					extendRight = true
					ihpSet = true
					// we are setup to create a new expression with left and right NUM operands and attach this node to the right of the existing node (expression)
					// this will be carried out during this current parse of NUM and the next NUM.
					operandL = true
				}
			}

			if operandL {

				numL = &num{i: i}
				fmt.Println("Left NUM", i, "  lvl: ", lvl)
				operandL = false

			} else {

				numR = &num{i: i}
				fmt.Println("Right NUM ", i, "  lvl: ", lvl)
				//
				en := makeExpr(lvl, numL, opr, numR)
				if e == nil {
					e, en = en, nil
				}
				// consumed following values, so reset them
				numL, numR = nil, nil
				opr = 0
				// AddParent is the default operation to extend the graph, which requires a numR only
				operandL = false
				// do not extend or add expression until we have an "e" and "en" expression.
				if en != nil {
					if extendRight {
						fmt.Printf("**** ExtendRight on  %c lvl %d     child: %c  %s lvl. %d \n\n", e.opr, e.depth, en.opr, en.name, en.depth)
						// higher precedence operator or a ( has occured - create new branch to right of current expression
						e, lvl = e.ExtendRight(en, lvl)
						// AddParent is the default method to extend the graph, so make extendRight false. Must be explicitly set to true
						// when the correct scenario occurs i.e. immediately after a ( or higher precedence operation detected
						extendRight = false
						//ihpSet = false //xxx
					} else {
						e = e.AddParent(en)
					}
				}
			}

		case token.PLUS:

			opr = PLUS
			if ihpSet {
				lvl--
				ihpSet = false
			}

		case token.MINUS:
			// is it a negative sign or a minus sign?
			ptok := prevTok.Type
			if ptok == token.LPAREN || ptok == token.MULTIPLY || ptok == token.PLUS || ptok == token.MINUS || ptok == token.DIVIDE {
				negative = true

			} else {

				opr = MINUS
				if ihpSet {
					lvl--
					ihpSet = false
				}
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
