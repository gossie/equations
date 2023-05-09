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
	subtraction := Sub(Num(4), Var(2, "x"))

	matcher := removeVariableSubtractionMatcher{}
	if !matcher.Match(&subtraction) {
		t.Fatal("matcher should match")
	}

	expected := Add(Num(4), Var(-2, "x"))
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
	division := Div(Num(4), Var(2, "x"))

	matcher := removeVariableDivisionMatcher{}
	if !matcher.Match(&division) {
		t.Fatal("matcher should match")
	}

	expected := Mul(Num(4), Var(1.0/2.0, "x"))
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

func TestReturnZeroMatcher_mulWith0_1(t *testing.T) {
	product := Mul(Add(Var(2, "x"), Num(4)), Num(0))

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
	product := Mul(Num(0), Add(Var(2, "x"), Num(4)))

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
	variable := Var(0, "x")

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

func TestReturnValueMatcher_mulWith1_1(t *testing.T) {
	product := Mul(Add(Var(2, "x"), Num(4)), Num(1))

	matcher := returnValueMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Add(Var(2, "x"), Num(4))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnValueMatcher_mulWith1_2(t *testing.T) {
	product := Mul(Num(1), Add(Var(2, "x"), Num(4)))

	matcher := returnValueMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Add(Var(2, "x"), Num(4))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnValueMatcher_add0_1(t *testing.T) {
	sum := Add(Add(Var(2, "x"), Num(4)), Num(0))

	matcher := returnValueMatcher{}
	if !matcher.Match(&sum) {
		t.Fatal("matcher should match")
	}

	expected := Add(Var(2, "x"), Num(4))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestReturnValueMatcher_add0_2(t *testing.T) {
	sum := Add(Num(0), Add(Var(2, "x"), Num(4)))

	matcher := returnValueMatcher{}
	if !matcher.Match(&sum) {
		t.Fatal("matcher should match")
	}

	expected := Add(Var(2, "x"), Num(4))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestVariableMulMatcher_1(t *testing.T) {
	product := Mul(Var(2, "x"), Num(3))

	matcher := variableMulMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Var(6, "x")
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestVariableMulMatcher_2(t *testing.T) {
	product := Mul(Num(3), Var(2, "x"))

	matcher := variableMulMatcher{}
	if !matcher.Match(&product) {
		t.Fatal("matcher should match")
	}

	expected := Var(6, "x")
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestVariableAddMatcher_1(t *testing.T) {
	sum := Add(Var(4, "x"), Var(2, "x"))

	matcher := variableAddMatcher{}
	if !matcher.Match(&sum) {
		t.Fatal("matcher should match")
	}

	expected := Var(6, "x")
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestDistributiveMatcher_1(t *testing.T) {
	formular := Mul(Add(Num(2), Num(4)), Num(3))

	matcher := distributiveMatcher{}
	if !matcher.Match(&formular) {
		t.Fatal("matcher should match")
	}

	expected := Add(Mul(Num(3), Num(2)), Mul(Num(3), Num(4)))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}

func TestDistributiveMatcher_2(t *testing.T) {
	formular := Mul(Num(3), Add(Num(2), Num(4)))

	matcher := distributiveMatcher{}
	if !matcher.Match(&formular) {
		t.Fatal("matcher should match")
	}

	expected := Add(Mul(Num(3), Num(2)), Mul(Num(3), Num(4)))
	result := matcher.Execute()
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expect %v to be %v", result, expected)
	}
}
