package equations

import "math"

type pattern func(*value) bool

func num(n float64) pattern {
	return func(v *value) bool {
		return v.op == "num" && v.number == n
	}
}

func anyNum(n *float64) pattern {
	return func(v *value) bool {
		if v.op == "num" {
			*n = v.number
			return true
		}
		return false
	}
}

func variable(name string) pattern {
	return func(val *value) bool {
		return val.op == "var" && val.name == name
	}
}

func anyVariable(factor *float64, name *string, exponent *float64) pattern {
	return func(v *value) bool {
		if v.op == "var" {
			*factor = v.number
			*name = v.name
			*exponent = v.exponent
			return true
		}
		return false
	}
}

func any(val *value) pattern {
	return func(v *value) bool {
		*val = *v
		return v != nil
	}
}

func bin(leftOperand pattern, operation string, rightOperand pattern) pattern {
	return func(val *value) bool {
		return val.op == operation && leftOperand(val.left) && rightOperand(val.right)
	}
}

type PatternMatcher interface {
	Match(*value) bool
	Execute() value
}

type removeSubtractionMatcher struct {
	valParam value
	number   float64
}

func (sm *removeSubtractionMatcher) Match(val *value) bool {
	return bin(any(&sm.valParam), "-", anyNum(&sm.number))(val)
}

func (sm *removeSubtractionMatcher) Execute() value {
	return Add(sm.valParam, Num(-sm.number))
}

type removeVariableSubtractionMatcher struct {
	valParam            value
	varFactor, exponent float64
	varName             string
}

func (sm *removeVariableSubtractionMatcher) Match(val *value) bool {
	return bin(any(&sm.valParam), "-", anyVariable(&sm.varFactor, &sm.varName, &sm.exponent))(val)
}

func (sm *removeVariableSubtractionMatcher) Execute() value {
	return Add(sm.valParam, Var(-sm.varFactor, sm.varName, sm.exponent))
}

type removeDivisionMatcher struct {
	valParam value
	number   float64
}

func (dm *removeDivisionMatcher) Match(val *value) bool {
	return bin(any(&dm.valParam), "/", anyNum(&dm.number))(val)
}

func (dm *removeDivisionMatcher) Execute() value {
	return Mul(dm.valParam, Num(1/dm.number))
}

type removeVariableDivisionMatcher struct {
	valParam            value
	varFactor, exponent float64
	varName             string
}

func (dm *removeVariableDivisionMatcher) Match(val *value) bool {
	return bin(any(&dm.valParam), "/", anyVariable(&dm.varFactor, &dm.varName, &dm.exponent))(val)
}

func (dm *removeVariableDivisionMatcher) Execute() value {
	return Mul(dm.valParam, Var(1/dm.varFactor, dm.varName, dm.exponent))
}

type addMatcher struct {
	number1, number2 float64
}

func (am *addMatcher) Match(val *value) bool {
	return bin(anyNum(&am.number1), "+", anyNum(&am.number2))(val)
}

func (am *addMatcher) Execute() value {
	return Num(am.number1 + am.number2)
}

type mulMatcher struct {
	number1, number2 float64
}

func (mm *mulMatcher) Match(val *value) bool {
	return bin(anyNum(&mm.number1), "*", anyNum(&mm.number2))(val)
}

func (mm *mulMatcher) Execute() value {
	return Num(mm.number1 * mm.number2)
}

type powMatcher struct {
	number1, number2 float64
}

func (pm *powMatcher) Match(val *value) bool {
	return bin(anyNum(&pm.number1), "^", anyNum(&pm.number2))(val)
}

func (pm *powMatcher) Execute() value {
	return Num(math.Pow(pm.number1, pm.number2))
}

type returnZeroMatcher struct {
}

func (mm *returnZeroMatcher) Match(val *value) bool {
	var val1, val2 value
	var number, exponent float64
	var varName string
	return bin(any(&val1), "*", num(0))(val) ||
		bin(num(0), "*", any(&val2))(val) ||
		(anyVariable(&number, &varName, &exponent)(val) && number == 0.0)
}

func (mm *returnZeroMatcher) Execute() value {
	return Num(0)
}

type returnOneMatcher struct {
}

func (mm *returnOneMatcher) Match(val *value) bool {
	var val1 value
	return bin(any(&val1), "^", num(0))(val)
}

func (mm *returnOneMatcher) Execute() value {
	return Num(1)
}

type returnValueMatcher struct {
	result value
}

func (rvm *returnValueMatcher) Match(val *value) bool {
	return bin(any(&rvm.result), "*", num(1))(val) ||
		bin(num(1), "*", any(&rvm.result))(val) ||
		bin(any(&rvm.result), "+", num(0))(val) ||
		bin(num(0), "+", any(&rvm.result))(val) ||
		bin(any(&rvm.result), "^", num(1))(val)
}

func (rvm *returnValueMatcher) Execute() value {
	return rvm.result
}

type variableMulMatcher struct {
	number1, number2, exponent float64
	varName                    string
}

func (mm *variableMulMatcher) Match(val *value) bool {
	return bin(anyVariable(&mm.number1, &mm.varName, &mm.exponent), "*", anyNum(&mm.number2))(val) ||
		bin(anyNum(&mm.number1), "*", anyVariable(&mm.number2, &mm.varName, &mm.exponent))(val)
}

func (mm *variableMulMatcher) Execute() value {
	return Var(mm.number1*mm.number2, mm.varName, mm.exponent)
}

type variableAddMatcher struct {
	number1, number2, exponent1, exponent2 float64
	varName1, varName2                     string
}

func (am *variableAddMatcher) Match(val *value) bool {
	return bin(anyVariable(&am.number1, &am.varName1, &am.exponent1), "+", anyVariable(&am.number2, &am.varName2, &am.exponent2))(val) && am.varName1 == am.varName2 && am.exponent1 == am.exponent2
}

func (am *variableAddMatcher) Execute() value {
	return Var(am.number1+am.number2, am.varName1, am.exponent1)
}

type variableMulVariableMatcher struct {
	factor1, exponent1, factor2, exponent2 float64
	varName1, varName2                     string
}

func (vmvm *variableMulVariableMatcher) Match(val *value) bool {
	return bin(anyVariable(&vmvm.factor1, &vmvm.varName1, &vmvm.exponent1), "*", anyVariable(&vmvm.factor2, &vmvm.varName2, &vmvm.exponent2))(val) && vmvm.varName1 == vmvm.varName2
}

func (vmvm *variableMulVariableMatcher) Execute() value {
	return Var(vmvm.factor1*vmvm.factor2, vmvm.varName1, vmvm.exponent1+vmvm.exponent2)
}

type distributiveMatcher struct {
	val1, val2 value
	number     float64
}

func (dm *distributiveMatcher) Match(val *value) bool {
	return bin(bin(any(&dm.val1), "+", any(&dm.val2)), "*", anyNum(&dm.number))(val) ||
		bin(anyNum(&dm.number), "*", bin(any(&dm.val1), "+", any(&dm.val2)))(val)
}

func (dm *distributiveMatcher) Execute() value {
	return Add(Mul(Num(dm.number), dm.val1), Mul(Num(dm.number), dm.val2))
}

type associativeMatcher1 struct {
	number1, number2, number3, exponent1, exponent2 float64
	varName1, varName2                              string
}

func (am *associativeMatcher1) Match(val *value) bool {
	return bin(bin(anyVariable(&am.number1, &am.varName1, &am.exponent1), "+", anyNum(&am.number2)), "+", anyVariable(&am.number3, &am.varName2, &am.exponent2))(val) && am.exponent1 == am.exponent2
}

func (am *associativeMatcher1) Execute() value {
	return Add(Var(am.number1+am.number3, am.varName1, am.exponent1), Num(am.number2))
}

type associativeMatcher2 struct {
	number1, number2 float64
	v                value
}

func (am *associativeMatcher2) Match(val *value) bool {
	return bin(bin(any(&am.v), "+", anyNum(&am.number1)), "+", anyNum(&am.number2))(val)
}

func (am *associativeMatcher2) Execute() value {
	return Add(am.v, Num(am.number1+am.number2))
}

type associativeMatcher3 struct {
	number1, number2 float64
	v                value
}

func (am *associativeMatcher3) Match(val *value) bool {
	return bin(bin(anyNum(&am.number1), "+", any(&am.v)), "+", anyNum(&am.number2))(val)
}

func (am *associativeMatcher3) Execute() value {
	return Add(Num(am.number1+am.number2), am.v)
}

type associativeMatcher4 struct {
	number1, number2, exponent1, exponent2 float64
	varName1, varName2                     string
	v                                      value
}

func (am *associativeMatcher4) Match(val *value) bool {
	return bin(bin(any(&am.v), "+", anyVariable(&am.number1, &am.varName1, &am.exponent1)), "+", anyVariable(&am.number2, &am.varName2, &am.exponent2))(val) && am.varName1 == am.varName2 && am.exponent1 == am.exponent2
}

func (am *associativeMatcher4) Execute() value {
	return Add(am.v, Var(am.number1+am.number2, am.varName1, am.exponent1))
}

type associativeMatcher5 struct {
	number1, number2 float64
	v                value
}

func (am *associativeMatcher5) Match(val *value) bool {
	return bin(anyNum(&am.number1), "+", bin(any(&am.v), "+", anyNum(&am.number2)))(val)
}

func (am *associativeMatcher5) Execute() value {
	return Add(am.v, Num(am.number1+am.number2))
}

type associativeMatcher6 struct {
	number1, number2 float64
	v                value
}

func (am *associativeMatcher6) Match(val *value) bool {
	return bin(anyNum(&am.number1), "+", bin(anyNum(&am.number2), "+", any(&am.v)))(val)
}

func (am *associativeMatcher6) Execute() value {
	return Add(am.v, Num(am.number1+am.number2))
}

var Matchers = []PatternMatcher{
	&removeSubtractionMatcher{},
	&removeVariableSubtractionMatcher{},
	&removeDivisionMatcher{},
	&removeVariableDivisionMatcher{},
	&addMatcher{},
	&mulMatcher{},
	&powMatcher{},
	&returnZeroMatcher{},
	&returnOneMatcher{},
	&returnValueMatcher{},
	&variableMulMatcher{},
	&variableAddMatcher{},
	// &variableAndNumberMulMatcher{},
	&distributiveMatcher{},
	&associativeMatcher1{},
	&associativeMatcher2{},
	&associativeMatcher3{},
	&associativeMatcher4{},
	&associativeMatcher5{},
	&associativeMatcher6{},
}
