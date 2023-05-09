package equations

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

func anyVariable(factor *float64, name *string) pattern {
	return func(v *value) bool {
		if v.op == "var" {
			*factor = v.number
			*name = v.name
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
	valParam  value
	varFactor float64
	varName   string
}

func (sm *removeVariableSubtractionMatcher) Match(val *value) bool {
	return bin(any(&sm.valParam), "-", anyVariable(&sm.varFactor, &sm.varName))(val)
}

func (sm *removeVariableSubtractionMatcher) Execute() value {
	return Add(sm.valParam, Var(-sm.varFactor, sm.varName))
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
	valParam  value
	varFactor float64
	varName   string
}

func (dm *removeVariableDivisionMatcher) Match(val *value) bool {
	return bin(any(&dm.valParam), "/", anyVariable(&dm.varFactor, &dm.varName))(val)
}

func (dm *removeVariableDivisionMatcher) Execute() value {
	return Mul(dm.valParam, Var(1/dm.varFactor, dm.varName))
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

type returnZeroMatcher struct {
}

func (mm *returnZeroMatcher) Match(val *value) bool {
	var val1, val2 value
	var number float64
	var varName string
	return bin(any(&val1), "*", num(0))(val) || bin(num(0), "*", any(&val2))(val) || (anyVariable(&number, &varName)(val) && number == 0.0)
}

func (mm *returnZeroMatcher) Execute() value {
	return Num(0)
}

type returnValueMatcher struct {
	result value
}

func (rvm *returnValueMatcher) Match(val *value) bool {
	return bin(any(&rvm.result), "*", num(1))(val) ||
		bin(num(1), "*", any(&rvm.result))(val) ||
		bin(any(&rvm.result), "+", num(0))(val) ||
		bin(num(0), "+", any(&rvm.result))(val)
}

func (rvm *returnValueMatcher) Execute() value {
	return rvm.result
}

type variableMulMatcher struct {
	number1, number2 float64
	varName          string
}

func (mm *variableMulMatcher) Match(val *value) bool {
	return bin(anyVariable(&mm.number1, &mm.varName), "*", anyNum(&mm.number2))(val) ||
		bin(anyNum(&mm.number1), "*", anyVariable(&mm.number2, &mm.varName))(val)
}

func (mm *variableMulMatcher) Execute() value {
	return Var(mm.number1*mm.number2, mm.varName)
}

type variableAddMatcher struct {
	number1, number2   float64
	varName1, varName2 string
}

func (am *variableAddMatcher) Match(val *value) bool {
	return bin(anyVariable(&am.number1, &am.varName1), "+", anyVariable(&am.number2, &am.varName2))(val) && am.varName1 == am.varName2
}

func (am *variableAddMatcher) Execute() value {
	return Var(am.number1+am.number2, am.varName1)
}

type variableAndNumberMulMatcher struct {
	number1, number2 float64
	varName          string
}

func (mm *variableAndNumberMulMatcher) Match(val *value) bool {
	return bin(anyVariable(&mm.number1, &mm.varName), "*", anyNum(&mm.number2))(val) ||
		bin(anyNum(&mm.number1), "*", anyVariable(&mm.number2, &mm.varName))(val)
}

func (mm *variableAndNumberMulMatcher) Execute() value {
	return Var(mm.number1*mm.number2, mm.varName)
}

type distributivMatcher struct {
	val1, val2 value
	number     float64
}

func (dm *distributivMatcher) Match(val *value) bool {
	return bin(bin(any(&dm.val1), "+", any(&dm.val2)), "*", anyNum(&dm.number))(val) ||
		bin(anyNum(&dm.number), "*", bin(any(&dm.val1), "+", any(&dm.val2)))(val)
}

func (dm *distributivMatcher) Execute() value {
	return Add(Mul(Num(dm.number), dm.val1), Mul(Num(dm.number), dm.val2))
}

type assoziativMatcher1 struct {
	number1, number2, number3 float64
	varName1, varName2        string
}

func (am *assoziativMatcher1) Match(val *value) bool {
	return bin(bin(anyVariable(&am.number1, &am.varName1), "+", anyNum(&am.number2)), "+", anyVariable(&am.number3, &am.varName2))(val)
}

func (am *assoziativMatcher1) Execute() value {
	return Add(Var(am.number1+am.number3, am.varName1), Num(am.number2))
}

type assoziativMatcher2 struct {
	number1, number2 float64
	v                value
}

func (am *assoziativMatcher2) Match(val *value) bool {
	return bin(bin(any(&am.v), "+", anyNum(&am.number1)), "+", anyNum(&am.number2))(val)
}

func (am *assoziativMatcher2) Execute() value {
	return Add(am.v, Num(am.number1+am.number2))
}

type assoziativMatcher3 struct {
	number1, number2 float64
	v                value
}

func (am *assoziativMatcher3) Match(val *value) bool {
	return bin(bin(anyNum(&am.number1), "+", any(&am.v)), "+", anyNum(&am.number2))(val)
}

func (am *assoziativMatcher3) Execute() value {
	return Add(Num(am.number1+am.number2), am.v)
}

type assoziativMatcher4 struct {
	number1, number2   float64
	varName1, varName2 string
	v                  value
}

func (am *assoziativMatcher4) Match(val *value) bool {
	return bin(bin(any(&am.v), "+", anyVariable(&am.number1, &am.varName1)), "+", anyVariable(&am.number2, &am.varName2))(val) && am.varName1 == am.varName2
}

func (am *assoziativMatcher4) Execute() value {
	return Add(am.v, Var(am.number1+am.number2, am.varName1))
}

type assoziativMatcher5 struct {
	number1, number2 float64
	v                value
}

func (am *assoziativMatcher5) Match(val *value) bool {
	return bin(anyNum(&am.number1), "+", bin(any(&am.v), "+", anyNum(&am.number2)))(val)
}

func (am *assoziativMatcher5) Execute() value {
	return Add(am.v, Num(am.number1+am.number2))
}

type assoziativMatcher6 struct {
	number1, number2 float64
	v                value
}

func (am *assoziativMatcher6) Match(val *value) bool {
	return bin(anyNum(&am.number1), "+", bin(anyNum(&am.number2), "+", any(&am.v)))(val)
}

func (am *assoziativMatcher6) Execute() value {
	return Add(am.v, Num(am.number1+am.number2))
}

var Matchers = []PatternMatcher{
	&removeSubtractionMatcher{},
	&removeVariableSubtractionMatcher{},
	&removeDivisionMatcher{},
	&removeVariableDivisionMatcher{},
	&addMatcher{},
	&mulMatcher{},
	&returnZeroMatcher{},
	&returnValueMatcher{},
	&variableMulMatcher{},
	&variableAddMatcher{},
	&variableAndNumberMulMatcher{},
	&distributivMatcher{},
	&assoziativMatcher1{},
	&assoziativMatcher2{},
	&assoziativMatcher3{},
	&assoziativMatcher4{},
	&assoziativMatcher5{},
	&assoziativMatcher6{},
}
