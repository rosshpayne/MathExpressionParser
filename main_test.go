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
		{input: "(2+2*7/2*10*(5+2))*2", result: 984},
		{input: "(2+2*7/2*10*(5*2))*2", result: 1404},
		{input: "(2+4*(5+2*3))*2", result: 92},
		{input: "(2+4*(5+2*3-5))*2", result: 52},
		{input: "(2+2*7/2*10*(5+2*3))*2", result: 1544},
		{input: "(2+2*7/2*10+6+4)*2", result: 164},
		{input: "(2+2*7/2*10+6+4*3)*2", result: 180},
		{input: "(((3*(7-3)*4)*2) + (5+2)*-8*2)*3", result: -48},
		{input: "(96 + (5+2)*-8*2)*3", result: -48},
		{input: "(96 + (5+2)*-8*2)+3", result: -13},
		{input: "((96 + (5+2))*-8*2)+3", result: -1645},
		{input: "(3*(7-3)*4)*2 + (5+2)*-8*2+3", result: -13},
		{input: "(5*7-3) + (4+3)*12", result: 116},
		{input: "(((3*7-3)*4)*2 + 5*2) +2 - -3", result: 159},
		{input: "(2+2*7/(2*10+6+4*3))*2", result: 4.7368},
		{input: "(2+2*7/((2*10+6+4*3)))*2", result: 4.7368},
		{input: "(2+2*7+(2*(10+6+4)-3))*2", result: 106},
		{input: "(2+2*7+(2*(10*6+4)-3))*2", result: 282},
		{input: "(2+2*7/(2*(2*6)-3))*2", result: 5.33333},
		{input: "(2+2*7/(2*2*6-3))*2", result: 5.33333},
		{input: "(2+2*7/(2*(2*6-3)))*2", result: 5.5555},
		{input: "(2+2*7/(2*((2*6)-3)))*2", result: 5.5555},
		{input: "(2+2*7/(2*(((2*6))-3)))*2", result: 5.5555},
		{input: "(2+2*7/(2*((2*6)-(3*5))))*2", result: -0.6666},
		{input: "(2+2*7/(2*((2*6+4-4)-(3*5))))*2", result: -0.6666},
		{input: "(2+2*7/(2*((2*6+4-4)-(3*5/5+10-20)+2*7)))*2", result: 4.42424},
		{input: "(2+2*7/(2*((2*6+4-4)-(2+3-5+3*5/5+10-20)+2*7)))*2", result: 4.42424},
	}
	for _, v := range tests {
		t.Log(v.input)
		root := numGraph(v.input)
		walk(root)
		if int(root.getResult()*10000) == int(v.result*10000) {
			t.Log("*** PASSED - ", v.result)
		} else {
			t.Errorf("+++FAILED - Got %g  expected %g\n", root.getResult(), v.result)
		}
	}
}

func TestNumAll(t *testing.T) {

	type testT struct {
		input  string
		result float64
	}

	var tests []testT = []testT{
		// {input: "5+4*3*3*2+3", result: 80},
		// {input: "(15-5)/(3-1)", result: 5},
		// {input: "2+4+(5+2)", result: 13},
		// {input: "2*4*(5+2)", result: 56},
		// {input: "2+4*(5+2)", result: 30},
		// {input: "2+4*(5-2)", result: 14},
		// {input: "2+(4*(5-2))", result: 14},
		// {input: "2+3*(4*(5-2))", result: 38},
		// {input: "5+4+3", result: 12},
		// {input: "5+4*3", result: 17},
		// {input: "5+4*3+3", result: 20},
		// {input: "5+4*3*4*2+3", result: 104},
		// {input: "3*(7-3)*4*2 ", result: 96},
		// {input: "((3*(7-3)))*2*3", result: 72},
		// {input: "((3+3*(7-3)))*2*3", result: 90},
		// {input: "((3*(7*2+3)))-2*32", result: -13},
		// {input: "(5+4*66/(3*2+4))-2*4", result: 23.4},
		// //{input: "5+4*66/(3*2+4)-2*4", result: 23.4},
		// {input: "3*(7-3)*2*3", result: 72},
		// {input: "(((3*7-3)*4)*2 + 5 )*2 - -3", result: 301},
		// {input: "(((3*7-3)*4)*2 + 5*2) - -3", result: 157},
		// {input: "2+5*7", result: 37},
		// {input: "2+5*7*2", result: 72},
		// {input: "2+2*7/2", result: 9},
		// {input: "2+2*7/2*10", result: 72},
		// {input: "(2+2*7/2*10)*2", result: 144},
		// {input: "(2+2*7/2*10*(5+2))*2", result: 984},
		// {input: "(2+2*7/2*10*(5*2))*2", result: 1404},
		// {input: "(2+4*(5+2*3))*2", result: 92},
		// {input: "(2+2*7/2*10*(5+2*3))*2", result: 1544},
		// {input: "(2+2*7/2*10+6+4)*2", result: 164},
		// {input: "(2+2*7/2*10+6+4*3)*2", result: 180},
		// {input: "(((3*(7-3)*4)*2) + (5+2)*-8*2)*3", result: -48},
		// {input: "(96 + (5+2)*-8*2)*3", result: -48},
		// {input: "(96 + (5+2)*-8*2)+3", result: -13},
		// {input: "((96 + (5+2))*-8*2)+3", result: -1645},
		// {input: "(3*(7-3)*4)*2 + (5+2)*-8*2+3", result: -13},
		// {input: "(5*7-3) + (4+3)*12", result: 116},
		// {input: "(((3*7-3)*4)*2 + 5*2) +2 - -3", result: 159},
		// {input: "(2+2*7/(2*10+6+4*3))*2", result: 4.7368},
		// {input: "(2+2*7/((2*10+6+4*3)))*2", result: 4.7368},
		// {input: "(2+2*7+(2*(10+6+4)-3))*2", result: 106},
		// {input: "(2+2*7+(2*(10*6+4)-3))*2", result: 282},
		// {input: "(2+2*7/(2*(2*6)-3))*2", result: 5.33333}, //Got -0.833333333333333
		// {input: "(2+2*7/(2*2*6-3))*2", result: 5.33333},   //Got -0.833333333333333
		// {input: "(2+2*7/(2*(2*6-3)))*2", result: 5.5555},
		//{input: "(2+2*7/(2*((2*6)-3)))*2", result: 5.5555}, //Got -0.833333333333333
	}
	for _, v := range tests {
		t.Log(v.input)
		root := numGraph(v.input)
		walk(root)
		if int(root.getResult()*10000) == int(v.result*10000) {
			t.Log("*** PASSED - ", v.result)
		} else {
			t.Errorf("+++FAILED - Got %g  expected %g\n", root.getResult(), v.result)
		}
	}
}
