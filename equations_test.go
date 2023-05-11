package equations_test

import (
	"fmt"
	"testing"

	"github.com/gossie/equations"
)

// func TestOptimze(t *testing.T) {
// 	left := equations.Add(equations.Var(4, "r"), equations.Mul(equations.Num(0), equations.Num(7)))
// 	right := equations.Add(equations.Var(1, "s"), equations.Div(equations.Num(25), equations.Num(5)))
// 	eq := equations.NewEquation(left, right)

// 	if eq.optimize().String() != "4.000000r = (1.000000s + 5.000000)" {
// 		t.Fatalf("expected %v to be 4.000000r = (1.000000s + 5.000000)", eq)
// 	}
// }

func TestSolveToR(t *testing.T) {
	left := equations.Add(equations.Var(4, "r", 1), equations.Mul(equations.Num(0), equations.Num(7)))
	right := equations.Add(equations.Var(1, "s", 1), equations.Div(equations.Num(25), equations.Num(5)))

	eq := equations.NewEquation(left, right)
	r, _ := equations.SolveTo(&eq, "r")

	if r.String() != "(0.250000s + 1.250000)" {
		t.Fatalf("expected %v to be (0.250000s + 1.250000)", r)
	}
}

func TestSolveToS(t *testing.T) {
	left := equations.Add(equations.Var(4, "r", 1), equations.Mul(equations.Num(0), equations.Num(7)))
	right := equations.Add(equations.Var(1, "s", 1), equations.Div(equations.Num(25), equations.Num(5)))

	original := equations.NewEquation(left, right)
	s, _ := equations.SolveTo(&original, "s")

	if s.String() != "(4.000000r + -5.000000)" {
		t.Fatalf("expected %v to be (4.000000r + -5.000000)", s)
	}

	if original.String() != "(4.000000r + (0.000000 * 7.000000)) = (1.000000s + (25.000000 / 5.000000))" {
		t.Fatalf("expected %v to be (4.000000r + (0.000000 * 7.000000)) = (1.000000s + (25.000000 / 5.000000))", original)
	}
}

func TestSolveTo_variableOnBothSides(t *testing.T) {
	left := equations.Add(equations.Var(4, "x", 1), equations.Mul(equations.Num(2), equations.Num(7)))
	right := equations.Add(equations.Var(2, "x", 1), equations.Div(equations.Num(25), equations.Num(5)))

	original := equations.NewEquation(left, right)
	s, _ := equations.SolveTo(&original, "x")

	if s.String() != fmt.Sprintf("%f", -9.0/2.0) {
		t.Fatalf("expected %v to be 4.5", s)
	}
}

func TestSet(t *testing.T) {
	left := equations.Add(equations.Var(4, "r", 1), equations.Mul(equations.Num(0), equations.Num(7)))
	right := equations.Add(equations.Var(1, "s", 1), equations.Div(equations.Num(25), equations.Num(5)))

	r := equations.Div(equations.Add(equations.Var(1, "s", 1), equations.Num(5)), equations.Num(4))

	eq := equations.NewEquation(left, right)
	eq = equations.Set(&eq, "r", r)
	if eq.String() != "((4.000000 * ((1.000000s + 5.000000) / 4.000000)) + (0.000000 * 7.000000)) = (1.000000s + (25.000000 / 5.000000))" {
		t.Fatalf("expected %v to be ((4.000000 * ((1.000000s + 5.000000) / 4.000000)) + (0.000000 * 7.000000)) = (1.000000s + (25.000000 / 5.000000))", eq)
	}

	// eq = eq.optimize()
	// if eq.String() != "(1.000000s + 5.000000) = (1.000000s + 5.000000)" {
	// 	t.Fatalf("expected %v to be (1.000000s + 5.000000) = (1.000000s + 5.000000)", eq)
	// }
}

// func TestOptimize_1(t *testing.T) {
// 	left := equations.Add(equations.Var(4, "r"), equations.Var(2, "r"))
// 	right := equations.Num(12.000000)

// 	eq := equations.NewEquation(left, right).optimize()
// 	if eq.String() != "6.000000r = 12.000000" {
// 		t.Fatalf("expected %v to be 6.000000r = 12.000000", eq)
// 	}
// }

// func TestOptimize_2(t *testing.T) {
// 	left := equations.Sub(equations.Var(4, "r"), equations.Var(2, "r"))
// 	right := equations.Num(12.000000)

// 	eq := equations.NewEquation(left, right).optimize()
// 	if eq.String() != "2.000000r = 12.000000" {
// 		t.Fatalf("expected %v to be 2.000000r = 12.000000", eq)
// 	}
// }

// func TestOptimize_3(t *testing.T) {
// 	left := equations.Mul(equations.Var(4, "r"), equations.Num(2))
// 	right := equations.Num(12.000000)

// 	eq := equations.NewEquation(left, right).optimize()
// 	if eq.String() != "8.000000r = 12.000000" {
// 		t.Fatalf("expected %v to be 8.000000r = 12.000000", eq)
// 	}
// }

// func TestOptimize_4(t *testing.T) {
// 	left := equations.Mul(equations.Num(2), equations.Var(4, "r"))
// 	right := equations.Num(12.000000)

// 	eq := equations.NewEquation(left, right).optimize()
// 	if eq.String() != "8.000000r = 12.000000" {
// 		t.Fatalf("expected %v to be 8.000000r = 12.000000", eq)
// 	}
// }

// func TestOptimize_5(t *testing.T) {
// 	left := equations.Div(equations.Var(4, "r"), equations.Num(2))
// 	right := equations.Num(12.000000)

// 	eq := equations.NewEquation(left, right).optimize()
// 	if eq.String() != "2.000000r = 12.000000" {
// 		t.Fatalf("expected %v to be 2.000000r = 12.000000", eq)
// 	}
// }

// func TestOptimize_6(t *testing.T) {
// 	left := equations.Mul(equations.Var(4, "r"), equations.Num(1))
// 	right := equations.Num(12.000000)

// 	eq := equations.NewEquation(left, right).optimize()
// 	if eq.String() != "4.000000r = 12.000000" {
// 		t.Fatalf("expected %v to be 4.000000r = 12.000000", eq)
// 	}
// }

// func TestOptimize_7(t *testing.T) {
// 	left := equations.Mul(equations.Num(1), equations.Var(4, "r"))
// 	right := equations.Num(12.000000)

// 	eq := equations.NewEquation(left, right).optimize()
// 	if eq.String() != "4.000000r = 12.000000" {
// 		t.Fatalf("expected %v to be 4.000000r = 12.000000", eq)
// 	}
// }

// func TestOptimize_8(t *testing.T) {
// 	left := equations.Mul(equations.Num(2), equations.Sub(equations.Var(1, "r"), equations.Num(4)))
// 	right := equations.Num(12.000000)

// 	eq := equations.NewEquation(left, right).optimize()
// 	if eq.String() != "(2.000000r + -8.000000) = 12.000000" {
// 		t.Fatalf("expected %v to be (2.000000r + -8.000000) = 12.000000", eq)
// 	}
// }

// func TestOptimize_9(t *testing.T) {
// 	left := equations.Mul(equations.Sub(equations.Var(1, "r"), equations.Num(4)), equations.Num(2))
// 	right := equations.Num(12.000000)

// 	eq := equations.NewEquation(left, right).optimize()
// 	if eq.String() != "(2.000000r + -8.000000) = 12.000000" {
// 		t.Fatalf("expected %v to be (2.000000r + -8.000000) = 12.000000", eq)
// 	}
// }
