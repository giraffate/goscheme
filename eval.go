package main

import (
	"errors"
)

func Eval(exp Exp, env *Env) (Exp, error) {
	switch t := exp.(type) {
	case Number, Boolean, Pair, string:
		return t, nil
	case Symbol:
		v := env.Get(t)
		if v == nil {
			return nil, errors.New("Variable is not found")
		}
		return v, nil
	case []Exp:
		if len(t) == 0 {
			return Pair{}, nil
		}

		switch t[0] {
		case Symbol("+"), Symbol("-"), Symbol("*"), Symbol("/"),
			Symbol("cons"), Symbol("car"), Symbol("cdr"), Symbol("list"),
			Symbol("null?"), Symbol("eq?"):
			return EvalPrimitive(t[0], t[1:], env)
		case Symbol("define"):
			return EvalDefine(t, env)
		case Symbol("if"):
			return EvalIf(t, env)
		case Symbol("cond"):
			return EvalCond(t, env)
		case Symbol("lambda"):
			return Lambda{t[1], t[2], env}, nil
		case Symbol("begin"):
			return EvalBegin(t, env)
		default:
			return EvalApplication(env, t[0], t[1:])
		}
	}
	return nil, nil
}

func EvalPrimitive(op Exp, args []Exp, env *Env) (Exp, error) {
	if args == nil || len(args) == 0 {
		return nil, errors.New("Invalid arguments")
	}

	var exps []Exp
	for _, arg := range args {
		exp, err := Eval(arg, env)
		if err != nil {
			return nil, err
		}
		exps = append(exps, exp)
	}

	switch op {
	case Symbol("+"):
		return applyPlus(exps)
	case Symbol("-"):
		return applyMinus(exps)
	case Symbol("*"):
		return applyMultiply(exps)
	case Symbol("/"):
		return applyDivision(exps)
	case Symbol("cons"):
		return applyCons(exps)
	case Symbol("car"):
		return applyCar(exps)
	case Symbol("cdr"):
		return applyCdr(exps)
	case Symbol("list"):
		return applyList(exps)
	case Symbol("null?"):
		return applyIsNull(exps)
	case Symbol("eq?"):
		return applyIsEqual(exps)
	default:
		return nil, errors.New("Invalid primitive instruction")
	}
}

func EvalBegin(exps []Exp, env *Env) (Exp, error) {
	if len(exps) < 2 {
		return nil, errors.New("Invalid arguments: begin")
	}

	beginExps := exps[1:]
	return EvalSequence(beginExps, env)
}

func EvalSequence(exps []Exp, env *Env) (Exp, error) {
	if len(exps) < 1 {
		return nil, errors.New("Invalid arguments")
	}

	if len(exps) == 1 {
		return Eval(exps[0], env)
	}

	if _, err := Eval(exps[0], env); err != nil {
		return nil, err
	}
	return EvalSequence(exps[1:], env)
}

func EvalCond(exps []Exp, env *Env) (Exp, error) {
	if len(exps) < 2 {
		return nil, errors.New("Invalid arguments: cond")
	}

	ifExp, err := makeCond(exps[1:], env)
	if err != nil {
		return nil, err
	}
	return Eval(ifExp, env)
}

func makeCond(exps []Exp, env *Env) (Exp, error) {
	first := exps[0]
	rest := exps[1:]

	firstExps, ok := first.([]Exp)
	if !ok {
		return nil, errors.New("Invalid type for cond")
	}

	if firstExps[0] == Symbol("else") {
		if len(rest) > 0 {
			return nil, errors.New("`else` is not last clause")
		}
		return makeBegin(firstExps[1:]), nil
	}

	elseExp, err := makeCond(rest, env)
	if err != nil {
		return nil, err
	}
	return makeIf(firstExps[0], makeBegin(firstExps[1:]), elseExp), nil
}

func makeIf(cond, thenExp, elseExp Exp) []Exp {
	return []Exp{Symbol("if"), cond, thenExp, elseExp}
}

func makeBegin(exps []Exp) []Exp {
	ret := []Exp{Symbol("begin")}
	for _, exp := range exps {
		ret = append(ret, exp)
	}
	return ret
}

// if-expression has (if cond then-exp else-exp) style
func EvalIf(exps []Exp, env *Env) (Exp, error) {
	if len(exps) != 4 {
		return nil, errors.New("Invalid arguments: if")
	}

	condExp, err := Eval(exps[1], env)
	if err != nil {
		return nil, err
	}

	var exp Exp
	if cond, ok := condExp.(Boolean); ok {
		if cond {
			exp = exps[2]
		} else {
			exp = exps[3]
		}
	} else {
		return nil, errors.New("Condition is not Boolean type")
	}

	return Eval(exp, env)
}

func EvalDefine(exps []Exp, env *Env) (Exp, error) {
	if len(exps) < 3 {
		return nil, errors.New("Define clause must have symbol and body")
	}
	t := exps[1]

	switch tt := t.(type) {
	// (define x 1) style
	case Symbol:
		v, err := Eval(exps[2], env)
		if err != nil {
			return nil, err
		}
		env.Put(t.(Symbol), v)
		return nil, nil

	// (define (doubled x) (* x x)) style
	case []Exp:
		name := tt[0]
		args := tt[1:]
		body := exps[2]
		env.Put(name.(Symbol), Lambda{Args: args, Body: body, Env: env})
		return nil, nil

	default:
		return nil, errors.New("Invalid Arguments: define")
	}
}

func EvalApplication(env *Env, operator Exp, operands []Exp) (Exp, error) {
	procedure, err := Eval(operator, env)
	if err != nil {
		return nil, err
	}

	var exps []Exp
	for _, operand := range operands {
		exp, err := Eval(operand, env)
		if err != nil {
			return nil, err
		}
		exps = append(exps, exp)
	}

	return Apply(procedure, exps)
}

func Apply(procedure Exp, args []Exp) (Exp, error) {
	if isCompound(procedure) {
		return compoundApply(procedure, args)
	}
	return nil, errors.New("Unknown procedure type")
}

func isCompound(procedure Exp) bool {
	switch procedure.(type) {
	case Lambda:
		return true
	}
	return false
}

func applyPlus(args []Exp) (Exp, error) {
	if len(args) < 2 {
		return nil, errors.New("Invalid arguments for +")
	}

	var sum int
	for _, arg := range args {
		n, ok := arg.(Number)
		if !ok {
			return nil, errors.New("+ only accepts number as arguments")
		}

		sum += int(n)
	}
	return Number(sum), nil
}

func applyMinus(args []Exp) (Exp, error) {
	if len(args) < 2 {
		return nil, errors.New("Invalid arguments for -")
	}

	var sum int
	ini, ok := args[0].(Number)
	if !ok {
		return nil, errors.New("- only accepts number as arguments")
	}
	sum = int(ini)

	for _, arg := range args[1:] {
		n, ok := arg.(Number)
		if !ok {
			return nil, errors.New("- only accepts number as arguments")
		}
		sum -= int(n)
	}
	return Number(sum), nil
}

func applyMultiply(args []Exp) (Exp, error) {
	if len(args) < 2 {
		return nil, errors.New("Invalid arguments for *")
	}

	var sum int
	ini, ok := args[0].(Number)
	if !ok {
		return nil, errors.New("* only accepts number as arguments")
	}
	sum = int(ini)

	for _, arg := range args[1:] {
		n, ok := arg.(Number)
		if !ok {
			return nil, errors.New("* only accepts number as arguments")
		}
		sum *= int(n)
	}
	return Number(sum), nil
}

func applyDivision(args []Exp) (Exp, error) {
	if len(args) < 2 {
		return nil, errors.New("Invalid arguments for /")
	}

	var sum int
	ini, ok := args[0].(Number)
	if !ok {
		return nil, errors.New("/ only accepts number as arguments")
	}
	sum = int(ini)

	for _, arg := range args[1:] {
		n, ok := arg.(Number)
		if !ok {
			return nil, errors.New("/ only accepts number as arguments")
		}
		sum /= int(n)
	}
	return Number(sum), nil
}

func applyCons(args []Exp) (Exp, error) {
	if len(args) != 2 {
		return nil, errors.New("Invalid arguments for cons")
	}
	return &Pair{Car: args[0], Cdr: args[1]}, nil
}

func applyCar(args []Exp) (Exp, error) {
	if len(args) != 1 {
		return nil, errors.New("Invalid arguments for car")
	}
	p, ok := args[0].(*Pair)
	if !ok {
		return nil, errors.New("car only accepts Pair as an argument")
	}
	return p.Car, nil
}

func applyCdr(args []Exp) (Exp, error) {
	if len(args) != 1 {
		return nil, errors.New("Invalid arguments for cdr")
	}
	p, ok := args[0].(*Pair)
	if !ok {
		return nil, errors.New("cdr only accepts Pair as an argument")
	}
	return p.Cdr, nil
}

func applyList(args []Exp) (Exp, error) {
	if len(args) < 1 {
		return nil, errors.New("Invalid arguments for list")
	}

	// reverse args
	tmp := make([]Exp, len(args), len(args))
	for i := 0; i <= len(args)/2; i++ {
		tmp[i], tmp[len(tmp)-1-i] = args[len(args)-1-i], args[i]
	}

	var list *Pair
	for _, v := range tmp {
		p := &Pair{v, list}
		list = p
	}
	return list, nil
}

func applyIsNull(args []Exp) (Exp, error) {
	if len(args) != 1 {
		return nil, errors.New("Invalid arguments for null?")
	}

	t, ok := args[0].(*Pair)
	if !ok {
		return nil, errors.New("null? only accepts Pair as an argument")
	}
	return Boolean(t.isNull()), nil
}

func applyIsEqual(args []Exp) (Exp, error) {
	if len(args) != 2 {
		return nil, errors.New("Invalid arguments for eq?")
	}
	return args[0] == args[1], nil
}

// procedure has (procedure params body env) format.
func compoundApply(procedure Exp, args []Exp) (Exp, error) {
	switch t := procedure.(type) {
	case Lambda:
		newEnv := ExpandEnv(t.Env)
		params := t.Args
		switch pt := params.(type) {
		case Symbol:
			if len(args) != 1 {
				return nil, errors.New("Invalid arguments")
			}
			newEnv.Put(pt, args[0])
		case []Exp:
			if len(pt) != len(args) {
				return nil, errors.New("Invalid arguments")
			}
			for _, k := range pt {
				for _, v := range args {
					newEnv.Put(k.(Symbol), v)
				}
			}
		}
		return Eval(t.Body, newEnv)
	default:
		return nil, errors.New("Unknown procedure type")
	}
}
