package equations

import (
	"reflect"
	"testing"
)

func TestRemoveSubtractionMatcher(t *testing.T) {
	subtraction := Sub(Num(4), Num(2))

	matcher := removeSubtractionMatcher{}
	if !matcher.Match(&subtraction) {
		t.Fatal("matcher should match")
	}

	expected := Add(Num(4), Num(-2))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestRemoveVariableSubtractionMatcher(t *testing.T) {
	subtraction := Sub(Num(4), Var(2, "x", 1))

	matcher := removeVariableSubtractionMatcher{}
	if !matcher.Match(&subtraction) {
		t.Fatal("matcher should match")
	}

	expected := Add(Num(4), Var(-2, "x", 1))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestRemoveDivisionMatcher(t *testing.T) {
	division := Div(Num(4), Num(2))

	matcher := removeDivisionMatcher{}
	if !matcher.Match(&division) {
		t.Fatal("matcher should match")
	}

	expected := Mul(Num(4), Num(1.0/2.0))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestRemoveVariableDivisionMatcher(t *testing.T) {
	division := Div(Num(4), Var(2, "x", 1))

	matcher := removeVariableDivisionMatcher{}
	if !matcher.Match(&division) {
		t.Fatal("matcher should match")
	}

	expected := Mul(Num(4), Var(1.0/2.0, "x", 1))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestAddMatcher(t *testing.T) {
	sum := Add(Num(4), Num(2))

	matcher := addMatcher{}
	if !matcher.Match(&sum) {
		t.Fatal("matcher should match")
	}

	expected := Num(6)
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestMulMatcher(t *testing.T) {
	product := Mul(Num(4), Num(2))

	matcher := mulMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Num(8)
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestPowMatcher(t *testing.T) {
	product := Pow(Num(2), Num(4))

	matcher := powMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Num(16)
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnZeroMatcher_mulWith0_1(t *testing.T) {
	product := Mul(Add(Var(2, "x", 1), Num(4)), Num(0))

	matcher := returnZeroMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Num(0)
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnZeroMatcher_mulWith0_2(t *testing.T) {
	product := Mul(Num(0), Add(Var(2, "x", 1), Num(4)))

	matcher := returnZeroMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Num(0)
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnZeroMatcher_variableFactorIsZero(t *testing.T) {
	variable := Var(0, "x", 1)

	matcher := returnZeroMatcher{}
	if !matcher.Match(&variable) {
		t.Fatal("matcher should match")
	}

	expected := Num(0)
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnOneMatcher(t *testing.T) {
	one := Pow(Add(Num(4), Num(2)), Num(0))

	matcher := returnOneMatcher{}
	if !matcher.Match(&one) {
		t.Fatal("matcher should match")
	}

	expected := Num(1)
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnValueMatcher_mulWith1_1(t *testing.T) {
	product := Mul(Add(Var(2, "x", 1), Num(4)), Num(1))

	matcher := returnValueMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Add(Var(2, "x", 1), Num(4))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnValueMatcher_mulWith1_2(t *testing.T) {
	product := Mul(Num(1), Add(Var(2, "x", 1), Num(4)))

	matcher := returnValueMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Add(Var(2, "x", 1), Num(4))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnValueMatcher_add0_1(t *testing.T) {
	sum := Add(Add(Var(2, "x", 1), Num(4)), Num(0))

	matcher := returnValueMatcher{}
	if !matcher.Match(&sum) {
		t.Fatal("matcher should match")
	}

	expected := Add(Var(2, "x", 1), Num(4))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnValueMatcher_add0_2(t *testing.T) {
	sum := Add(Num(0), Add(Var(2, "x", 1), Num(4)))

	matcher := returnValueMatcher{}
	if !matcher.Match(&sum) {
		t.Fatal("matcher should match")
	}

	expected := Add(Var(2, "x", 1), Num(4))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnValueMatcher_powWith1(t *testing.T) {
	product := Pow(Add(Var(2, "x", 1), Num(4)), Num(1))

	matcher := returnValueMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Add(Var(2, "x", 1), Num(4))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestVariableMulMatcher_1(t *testing.T) {
	product := Mul(Var(2, "x", 1), Num(3))

	matcher := variableMulMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Var(6, "x", 1)
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestVariableMulMatcher_2(t *testing.T) {
	product := Mul(Num(3), Var(2, "x", 1))

	matcher := variableMulMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Var(6, "x", 1)
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestVariableMulVariableMatcher_3(t *testing.T) {
	product := Mul(Var(3, "x", 2), Var(2, "x", 1))

	matcher := variableMulVariableMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Var(6, "x", 3)
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestVariableMulVariableMatcher_4(t *testing.T) {
	product := Mul(Var(3, "x", 2), Var(2, "y", 1))

	matcher := variableMulVariableMatcher{}
	if matcher.Match(&product) {
		t.Fatal("matcher should not match")
	}
}

func TestVariableAddMatcher_1(t *testing.T) {
	sum := Add(Var(4, "x", 1), Var(2, "x", 1))

	matcher := variableAddMatcher{}
	if !matcher.Match(&sum) {
		t.Fatal("matcher should match")
	}

	expected := Var(6, "x", 1)
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestVariableAddMatcher_2(t *testing.T) {
	sum := Add(Var(4, "x", 1), Var(2, "x", 2))

	matcher := variableAddMatcher{}
	if matcher.Match(&sum) {
		t.Fatal("matcher should not match")
	}
}

func TestDistributiveMatcher_1(t *testing.T) {
	formula := Mul(Add(Num(2), Num(4)), Num(3))

	matcher := distributiveMatcher{}
	if !matcher.Match(&formula) {
		t.Fatal("matcher should match")
	}

	expected := Add(Mul(Num(3), Num(2)), Mul(Num(3), Num(4)))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestDistributiveMatcher_2(t *testing.T) {
	formula := Mul(Num(3), Add(Num(2), Num(4)))

	matcher := distributiveMatcher{}
	if !matcher.Match(&formula) {
		t.Fatal("matcher should match")
	}

	expected := Add(Mul(Num(3), Num(2)), Mul(Num(3), Num(4)))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestBinomial1(t *testing.T) {
	formula := Pow(Add(Var(2, "x", 1), Num(3)), Num(2))

	matcher := binomial1Matcher{}
	if !matcher.Match(&formula) {
		t.Fatal("matcher should match")
	}

	expected := Add(Add(Pow(Var(2, "x", 1), Num(2)), Mul(Num(2), Mul(Var(2, "x", 1), Num(3)))), Pow(Num(3), Num(2)))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestBinomial3_1(t *testing.T) {
	formula := Mul(Add(Var(2, "x", 1), Num(3)), Add(Var(2, "x", 1), Num(-3)))

	matcher := binomial3Matcher{}
	if !matcher.Match(&formula) {
		t.Fatal("matcher should match")
	}

	expected := Sub(Pow(Var(2, "x", 1), Num(2)), Pow(Num(3), Num(2)))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestBinomial3_2(t *testing.T) {
	formula := Mul(Add(Var(2, "x", 1), Var(1, "y", 1)), Sub(Var(2, "x", 1), Var(1, "y", 1)))

	matcher := binomial3Matcher{}
	if !matcher.Match(&formula) {
		t.Fatal("matcher should match")
	}

	expected := Sub(Pow(Var(2, "x", 1), Num(2)), Pow(Var(1, "y", 1), Num(2)))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}
