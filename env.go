package eval

import (
	"fmt"

	c "github.com/nichtsen/lis/vcons"
)

var GlobalEnv *Environment

func init() {
	InitGlobal()
}

func InitGlobal() {
	GlobalEnv = &Environment{
		Frame: Frame{
			"+":         Procedure(add),
			"-":         Procedure(substract),
			"*":         Procedure(multiply),
			"==":        Procedure(equal),
			">":         Procedure(larger),
			"append":    Procedure(sappend),
			"cons":      Procedure(cons),
			"car":       Procedure(car),
			"cdr":       Procedure(cdr),
			"list":      Procedure(list),
			"null?":     Procedure(isNull),
			"not-null?": Procedure(isNotNull),
			"print":     Procedure(print),
		},
		Enclose: nil,
	}
}

func add(args ...interface{}) interface{} {
	return args[0].(int) + args[1].(int)
}

func sappend(args ...interface{}) interface{} {
	return args[0].(string) + args[1].(string)
}

func equal(args ...interface{}) interface{} {
	return args[0].(int) == args[1].(int)
}

func larger(args ...interface{}) interface{} {
	return args[0].(int) > args[1].(int)
}

func multiply(args ...interface{}) interface{} {
	return args[0].(int) * args[1].(int)
}

func substract(args ...interface{}) interface{} {
	return args[0].(int) - args[1].(int)
}

func cons(args ...interface{}) interface{} {
	return c.Cons(args[0], args[1])
}

func car(args ...interface{}) interface{} {
	return c.Car(args[0])
}

func cdr(args ...interface{}) interface{} {
	return c.Cdr(args[0])
}

func list(args ...interface{}) interface{} {
	return c.List(args...)
}

func isNull(args ...interface{}) interface{} {
	return c.Empty(args[0])
}

func isNotNull(args ...interface{}) interface{} {
	return !c.Empty(args[0])
}

func print(args ...interface{}) interface{} {
	for _, arg := range args {
		fmt.Print(arg)
	}
	return nil
}

type Frame map[string]interface{}
type Environment struct {
	Enclose *Environment
	Frame   Frame
}

func NewFrame(vars []string, vals []interface{}) Frame {
	if len(vars) != len(vals) {
		panic("length of vars and vals should be equal")
	}
	res := make(Frame)
	for idx, v := range vars {
		res[v] = vals[idx]
	}
	return res
}

func ExtendEnv(enclose *Environment, vars []string, vals []interface{}) *Environment {
	return &Environment{
		Enclose: enclose,
		Frame:   NewFrame(vars, vals),
	}
}

func IsEmpty(env *Environment) bool {
	return env == nil
}

func (env *Environment) LookUpVariable(v string) interface{} {
	if val, ok := env.Frame[v]; ok {
		return val
	} else {
		if IsEmpty(env.Enclose) {
			panic(fmt.Sprintf("Undefined varible %v", v))
		} else {
			return env.Enclose.LookUpVariable(v)
		}
	}
}

func (env *Environment) SetVariable(v string, val interface{}) {
	if _, ok := env.Frame[v]; ok {
		env.Frame[v] = val
	} else {
		if IsEmpty(env.Enclose) {
			panic(fmt.Sprintf("Unbound varible %v", v))
		} else {
			env.Enclose.SetVariable(v, val)
		}
	}
}

func (env *Environment) DefineVariable(v string, val interface{}) {
	if _, ok := env.Frame[v]; ok {
		panic(fmt.Sprintf("varible has been defined %v", v))
	}
	env.Frame[v] = val
}
