package eval

import (
	"fmt"
	"strconv"
)

func Eval(e *Expression, env *Environment) interface{} {
	for expr := (*e); len(expr) > 0; expr = (*e) {
		switch {
		case DefineExpr(expr):
			EvalDefine(e, env)
		case AssignExpr(expr):
			EvalAssign(e, env)
		case LambdaExpr(expr):
			val := EvalLambda(e, env)
			exprNext := (*e)
			if LambdaApplicationExpr(exprNext) {
				paras, idx := ApplicationaParas(exprNext)
				*e = (*e)[idx:]
				return Apply(val, EvalArgs(paras, env)...)
			}
			return val
		case IfExpr(expr):
			val, ifContinue := EvalIf(e, env)
			if !ifContinue {
				return val
			}
		case NumberExpr(expr):
			n, _ := strconv.Atoi(expr[0])
			*e = expr[1:]
			return n
		case SymbolExpr(expr):
			*e = expr[1:]
			return expr[0][1:]
		case ApplicationExpr(expr):
			ae := ApplicationName(expr)
			*e = (*e)[1:]
			paras, idx := ApplicationaParas(*e)
			*e = (*e)[idx:]
			return Apply(Eval(&ae, env), EvalArgs(paras, env)...)
		case PerformExpr(expr):
			expr = expr[1:]
			ae := ApplicationName(expr)
			*e = (*e)[2:]
			paras, idx := ApplicationaParas(*e)
			*e = (*e)[idx:]
			Apply(Eval(&ae, env), EvalArgs(paras, env)...)
		default:
			val := env.LookUpVariable(expr[0])
			*e = expr[1:]
			return val
		}
	}
	return nil
}

func EvalArgs(exprs []Expression, env *Environment) []interface{} {
	res := make([]interface{}, 0)
	for _, expr := range exprs {
		res = append(res, Eval(&expr, env))
	}
	return res
}

func Apply(app interface{}, args ...interface{}) interface{} {
	if IsPrimitive(app) {
		return Primitive(app)(args...)
	}
	if IsCompound(app) {
		cp := Compound(app)
		e := cp.Body()
		return Eval(&e, ExtendEnv(cp.Env(), cp.Paras(), args))
	}
	panic(fmt.Sprintf("Invalid application %v", app))
}

func makeLambda(e *Expression, env *Environment) CompoundProcedure {
	paras, idx := ProcedureParas(*e)
	*e = (*e)[idx:]
	body, idx := ProcedureBody(*e)
	cp := NewCompoundProdedure(paras, body, env)
	*e = (*e)[idx:]
	return cp
}

func EvalDefine(e *Expression, env *Environment) {
	(*e) = (*e)[1:]
	// define a procedure variable
	if ProcedureExpr((*e)) {
		va := ProcedureVar(*e)
		cp := makeLambda(e, env)
		env.DefineVariable(va, cp)
		return
	}
	// define a non-proc variable
	va := (*e)[0]
	*e = (*e)[1:]
	env.DefineVariable(va, Eval(e, env))
}

func EvalLambda(e *Expression, env *Environment) CompoundProcedure {
	return makeLambda(e, env)
}

func EvalAssign(e *Expression, env *Environment) {
	va := (*e)[1]
	*e = (*e)[2:]
	env.SetVariable(va, Eval(e, env))
}

func EvalIf(e *Expression, env *Environment) (interface{}, bool) {
	*e = (*e)[1:]
	val := Eval(e, env)
	v, ok := val.(bool)
	if !ok {
		panic(fmt.Sprintf("Invalid Prediction, %v", v))
	}
	csq, idx := IfBody(*e)
	defer func() {
		*e = (*e)[idx:]
	}()
	if v {
		val := Eval(&csq, env)
		if val != nil {
			return val, false
		}
	}
	return nil, true
}
