package equations

import (
	"errors"
	"fmt"
)

type BinaryOp func(value, value) value

var operators = map[string]BinaryOp{
	"+": Add,
	"*": Mul,
	"-": Sub,
	"/": Div,
}

var rightComplements = map[string]string{
	"+": "-",
	"*": "/",
	"-": "+",
	"/": "*",
}

var leftComplements = map[string]string{
	"+": "-",
	"*": "/",
	"-": "-",
	"/": "/",
}

type opValuePair struct {
	op   string
	val  value
	swap bool
}

type path []*opValuePair

func findValue(val *value, name string) (*value, path, path, error) {
	if variable(name)(val) {
		return val, make(path, 0), append(make(path, 0), &opValuePair{"/", Num(val.number), false}), nil
	}

	if val.left != nil || val.right != nil {
		if val.left != nil {
			found, pathToValue, complementaryPath, err := findValue(val.left, name)
			if err == nil {
				op := rightComplements[val.op]
				pathToValue = append(pathToValue, &opValuePair{op, *val.left, val.op == "-" || val.op == "/"})
				complementaryPath = append(complementaryPath, &opValuePair{op, *val.right, false})
				return found, pathToValue, complementaryPath, nil
			}
		}

		if val.right != nil {
			found, pathToValue, complementaryPath, err := findValue(val.right, name)
			if err == nil {
				op := leftComplements[val.op]
				pathToValue = append(pathToValue, &opValuePair{op, *val.right, false})
				complementaryPath = append(complementaryPath, &opValuePair{op, *val.left, val.op == "-" || val.op == "/"})
				return found, pathToValue, complementaryPath, nil
			}
		}
	}

	return nil, nil, nil, errors.New("variable " + name + " not found")
}

type SolveError struct {
	err           error
	FinalEquation equation
}

func (se *SolveError) Error() string {
	return se.err.Error()
}

type equation struct {
	left, right value
}

func NewEquation(left, right value) equation {
	return equation{left: left, right: right}
}

func (e equation) IsTrue() bool {
	return e.left.execute() == e.right.execute()
}

func (e equation) optimize() equation {
	l := e.left.execute()
	r := e.right.execute()
	return NewEquation(l, r)
}

func (eq equation) SolveTo(varName string) (*value, error) {
	left, _, leftComplementaryPath, errLeft := findValue(&eq.left, varName)
	right, rightPath, rightComplementaryPath, errRight := findValue(&eq.right, varName)

	if left != nil && right != nil {
		eq := NewEquation(processPathElement(rightPath[len(rightPath)-1], eq.left), processPathElement(rightPath[len(rightPath)-1], eq.right))
		eq = eq.optimize()
		return eq.SolveTo(varName)
	}

	if errLeft != nil && errRight != nil {
		return nil, &SolveError{errors.New(varName + " could not be found"), eq}
	}

	if left != nil {
		return processPath(eq.right, leftComplementaryPath), nil
	} else {
		return processPath(eq.left, rightComplementaryPath), nil
	}
}

func (e equation) Set(varName string, val value) equation {
	newLeft := insert(e.left, varName, val)
	newRight := insert(e.right, varName, val)
	return NewEquation(newLeft, newRight)
}

func insert(current value, varName string, val value) value {
	if variable(varName)(&current) {
		return Mul(Num(current.number), val)
	}

	if op, present := operators[current.op]; present {
		return op(insert(*current.left, varName, val), insert(*current.right, varName, val))
	}

	return current
}

func processPath(val value, p path) *value {
	current := val
	for i := len(p) - 1; i >= 0; i-- {
		current = processPathElement(p[i], current)
	}
	result := current.execute()
	return &result
}

func processPathElement(v *opValuePair, current value) value {
	switch v.op {
	default:
		panic("unkown operator " + v.op)
	case "+":
		if v.swap {
			return Add(v.val, current)
		} else {
			return Add(current, v.val)
		}
	case "*":
		if v.swap {
			return Mul(v.val, current)
		} else {
			return Mul(current, v.val)
		}
	case "-":
		if v.swap {
			return Sub(v.val, current)
		} else {
			return Sub(current, v.val)
		}
	case "/":
		if v.swap {
			return Div(v.val, current)
		} else {
			return Div(current, v.val)
		}
	}
}

func (e equation) String() string {
	return fmt.Sprintf("%v = %v", e.left, e.right)
}

type value struct {
	left, right *value
	op          string
	number      float64
	name        string
}

func (v value) Number() float64 {
	switch v.op {
	default:
		panic("value " + v.String() + " is not terminal")
	case "num":
		return v.number
	}
}

func (v value) execute() value {
	if v.left != nil && v.right != nil {
		l := v.left.execute()
		r := v.right.execute()
		v.left = &l
		v.right = &r
	}

	for _, pm := range Matchers {
		if pm.Match(&v) {
			return pm.Execute().execute()
		}
	}
	return v
}

func (v value) String() string {
	switch v.op {
	default:
		panic("unknown operator: " + v.op)
	case "num":
		return fmt.Sprintf("%f", v.number)
	case "var":
		return fmt.Sprintf("%f%v", v.number, v.name)
	case "+":
		return fmt.Sprintf("(%v + %v)", v.left, v.right)
	case "*":
		return fmt.Sprintf("(%v * %v)", v.left, v.right)
	case "-":
		return fmt.Sprintf("(%v - %v)", v.left, v.right)
	case "/":
		return fmt.Sprintf("(%v / %v)", v.left, v.right)
	}
}

func Num(number float64) value {
	return value{number: number, op: "num"}
}

func Add(left, right value) value {
	return value{left: &left, right: &right, op: "+"}
}

func Sub(left, right value) value {
	return value{left: &left, right: &right, op: "-"}
}

func Mul(left, right value) value {
	return value{left: &left, right: &right, op: "*"}
}

func Div(left, right value) value {
	return value{left: &left, right: &right, op: "/"}
}

func Var(factor float64, name string) value {
	return value{number: factor, name: name, op: "var"}
}
