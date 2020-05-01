package main

import (
	"testing"
)

// Numeric Expressions
func TestNumX(t *testing.T) {

	type testT struct {
		input  string
		result float64
	}

	var tests []testT = []testT{
		{input: "15/(3+2)*3", result: 9},
		{input: "15/(3+2)-3", result: 0},
		{input: "15/((3+2)-3)", result: 7.5},
		{input: "2*((2*6)-3)", result: 18},
		{input: "(2*((2*6)-3))/2", result: 9},
		{input: "2+2+15/(3+2)*3+4", result: 17},
		{input: "2+2*15/(3+2)*3+4", result: 24},
		{input: "2+(2*15/(3+2)*3)+4", result: 24},
		{input: "15/(3+2)+3", result: 6},
		{input: "15/(3+2)+3*2", result: 9},
		{input: "15/(3+2)*2+3*2", result: 12},
		{input: "2+((3+15/(3+2))*2*3)*2+3*2", result: 80},
		{input: "2+(((3+15/(3+2))*2)*3)*2+3*2", result: 80},
		{input: "2+(((3+15/(3+2))*2*3))*2+3*2", result: 80},
		{input: "2+(3+15/(3+2))*2*3*2+3*2", result: 80},
		{input: "15+(4+2)*3", result: 33},
		{input: "15+(4+2)*3+2+5*2", result: 45},
		{input: "15+(4+2*3+2)*3", result: 51},
		{input: "15+3*(4+2*3+2)", result: 51},
		{input: "5+4*3*3*2+3", result: 80},
		{input: "(15-5)/(3-1)", result: 5},
		{input: "2+4+(5+2)", result: 13},
		{input: "2*4*(5+2)", result: 56},
		{input: "2+4*(5+2)", result: 30},
		{input: "2+4*(5-2)", result: 14},
		{input: "2+(4*(5-2))", result: 14},
		{input: "2+3*(4*(5-2))", result: 38},
		{input: "5+4+3", result: 12},
		{input: "5+4*3", result: 17},
		{input: "5+4*3+3", result: 20},
		{input: "5+4*3*4*2+3", result: 104},
		{input: "3*(7-3)*4*2 ", result: 96},
		{input: "((3*(7-3)))*2*3", result: 72},
		{input: "((3+3*(7-3)))*2*3", result: 90},
		{input: "((3*(7*2+3)))-2*32", result: -13},
		{input: "(5+4*66/(3*2+4))-2*4", result: 23.4},
		{input: "5+4*66/(3*2+4)-2*4", result: 23.4},
		{input: "3*(7-3)*2*3", result: 72},
		{input: "(((3*7-3)*4)*2 + 5 )*2 - -3", result: 301},
		{input: "(((3*7-3)*4)*2 + 5*2) - -3", result: 157},
		{input: "2+5*7", result: 37},
		{input: "2+5*7*2", result: 72},
		{input: "2+2*7/2", result: 9},
		{input: "2+2*7/2*10", result: 72},
		{input: "(2+2*7/2*10)*2", result: 144},
		{input: "(2+2*7/2*10)+2-6*2", result: 62},
		{input: "(2+2*7/(2*10))+2-6*2", result: -7.3},
		{input: "(2+2*7/(2*10+1-1+2-2)*2+1+3)+2-6*2", result: -2.59999},
		{input: "(2+2*7/(2*10)*2+1+3)+2-6*2", result: -2.59999},
		{input: "(2+2*7/(2*10)*2+1+3)+2-6*2", result: -2.59999},
		{input: "(2+2*7+(2*10+2)*2+1+3)+2-6*2", result: 54},
		{input: "(2+2*7/2*10*(5+2))*2", result: 984},
		{input: "(2+2*7/2*10*(5*2))*2", result: 1404},
		{input: "(2+4*(5+2*3))*2", result: 92},
		{input: "(2+4*(5+2*3-5))*2", result: 52},
		{input: "(2+2*7/2*10*(5+2*3))*2", result: 1544},
		{input: "(2+2*7/2*10+6+4)*2", result: 164},
		{input: "(2+2*7/2*10+6+4*3)*2", result: 180},
		{input: "(3*(7-3)*4)*2", result: 96},
		{input: "3*(7-3)*4*2 + (5+2)*-8", result: 40},
		{input: "(3*(7-3)*4)*2 + (5+2)*-8", result: 40},
		{input: "((3*(7-3)*4)*2)+ (5+2)*-8", result: 40},
		{input: "(3*(7-3)*4*2) + (5+2)*-8", result: 40},
		{input: "(3*(7-3)*4*2) + (5+2)*-8*2", result: -16},
		{input: "((3*(7-3)*4*2) + (5+2)*-8*2)", result: -16},
		{input: "((3*(7-3)*4)*2) + (5+2)*-8*2", result: -16},
		{input: "(((3*(7-3)*4)*2) + (5+2)*-8*2)", result: -16},
		{input: "(((3*(7-3)*4)*2) + (5+2)*-8*2)*3", result: -48},
		{input: "(((3*(7-3*2)*4)*2) + (5+2)*-8*2)*3", result: -264},
		{input: "((3*(7-3*2)*4*2) + (5+2)*-8*2)*3", result: -264},
		{input: "(3*(7-3*2)*4*2 + (5+2)*-8*2)*3", result: -264},
		{input: "(3+(19-3*2)*4*2 + (5+2)*-8*2)+66", result: 61},
		{input: "(((3*(7-3*2+1-1)*4)*2) + (5+2)*-8*2)*3", result: -264},
		{input: "(((3*(7-3)+4)*2) + (5+2)*-8*2)*3", result: -240},
		{input: "(96 + (5+2)*-8*2)*3", result: -48},
		{input: "(96 + (5+2)*-8*2)+3", result: -13},
		{input: "((96 + (5+2))*-8*2)+3", result: -1645},
		{input: "(3*(7-3)*4)*2 + (5+2)*-8*2+3", result: -13},
		{input: "(5*7-3) + (4+3)*12", result: 116},
		{input: "(((3*7-3)*4)*2 + 5*2) +2 - -3", result: 159},
		{input: "(2+2*7/(2*10+6+4*3))*2", result: 4.7368},
		{input: "(2+2*7+(2*(10+6+4)-3))*2", result: 106},
		{input: "(2+2*7/((2*10+6+4*3)))*2", result: 4.7368},
	}
	for _, v := range tests {
		t.Log(v.input)
		root := buildExprGraph(v.input)
		walk(root)
		if int(root.getResult()*10000) == int(v.result*10000) {
			t.Log("*** PASSED - ", v.result)
		} else {
			t.Errorf("+++FAILED - Got %g  expected %g\n", root.getResult(), v.result)
		}
	}
}
