package eval

import (
	"fmt"
	"strings"
	"unicode"
)

type Expression []string

func MakeExpr(text string) *Expression {
	val := Expression(
		strings.Fields(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(
						strings.ReplaceAll(
							strings.ReplaceAll(
								strings.ReplaceAll(text, ";", " "),
								"{", " { ",
							),
							"}", " } ",
						),
						"(", " ( ",
					),
					")", " ) ",
				),
				",", " , ",
			),
		),
	)
	return &val
}

func DefineExpr(expr Expression) bool { return expr[0] == "define" }
func AssignExpr(expr Expression) bool { return expr[0] == "set" }
func LambdaExpr(expr Expression) bool {
	return expr[0] == "lambda" && ProcedureExpr(expr)
}
func IfExpr(expr Expression) bool     { return expr[0] == "if" }
func NumberExpr(expr Expression) bool { return strings.IndexFunc(expr[0], unicode.IsNumber) == 0 }
func SymbolExpr(expr Expression) bool { return strings.IndexRune(expr[0], '\'') == 0 }
func ApplicationExpr(expr Expression) bool {
	if len(expr) < 3 {
		return false
	}
	return expr[1] == "("
}
func LambdaApplicationExpr(expr Expression) bool {
	if len(expr) < 2 {
		return false
	}
	return expr[0] == "("
}

func ApplicationName(expr Expression) Expression {
	return Expression{expr[0]}
}

func ApplicationaParas(expr Expression) ([]Expression, int) {
	subExpr, idx := ExtractParas(expr)
	return complexParas(subExpr), idx + 1
}

func complexParas(expr Expression) []Expression {
	res := make([]Expression, 0)
	loopParas(expr, &res)
	return res
}

func simpleParas(expr Expression) []string {
	var b strings.Builder
	for _, str := range expr {
		b.WriteString(str)
	}
	return strings.Split(b.String(), ",")
}

func loopParas(seg Expression, e *[]Expression) {
	if len(seg) == 0 {
		return
	}
	var count int
	for idx, str := range seg {
		if str == "(" {
			count++
		}
		if str == ")" {
			count--
		}
		// a comma not in any parenthese scope
		if str == "," && count == 0 {
			*e = append(*e, seg[:idx])
			loopParas(seg[idx+1:], e)
			return
		}
	}
	if count != 0 {
		panic(fmt.Sprintf("unenclosed bracket:\n %s", seg))
	}
	*e = append(*e, seg)
}

// deprecated
func innerParas(strs []string) []string {
	res := make([]string, 0)
	var count int
	var scan int
	for idx, str := range strs {
		if count > 0 {
			res[scan] = res[scan] + "," + strs[idx]
		}
		if count == 0 {
			res = append(res, strs[idx])
		}
		if count < 0 {
			panic(fmt.Sprintf("unenclosed parenthese: %v", strs))
		}
		if strings.Index(str, "(") > 0 {
			if count == 0 {
				scan = len(res) - 1
			}
			count++
		}
		if strings.Index(str, ")") > 0 {
			count--
		}
	}
	return res
}

func ProcedureExpr(expr Expression) bool {
	if len(expr) < 5 {
		return false
	}
	if expr[1] != "(" {
		return false
	}
	_, idx := ExtractParas(expr[1:])
	if len(expr) < (idx+1+1)+1 {
		return false
	}
	return expr[idx+1+1] == "{"
}

func ProcedureVar(expr Expression) string {
	return expr[0]
}

func ProcedureParas(expr Expression) ([]string, int) {
	subExpr, idx := ExtractParas(expr[1:])
	return simpleParas(subExpr), idx + 2
}

func ProcedureBody(expr Expression) (Expression, int) {
	val, idx := NextBlock(expr)
	return val, idx + 1
}

func PerformExpr(expr Expression) bool {
	return expr[0] == "perform" && len(expr) > 1 && ApplicationExpr(expr[1:])
}

func IfBody(expr Expression) (Expression, int) {
	val, idx := NextBlock(expr)
	return val, idx + 1
}

func (e Expression) Rest() Expression {
	return e[1:]
}

func (e Expression) Advance(n int) Expression {
	if n > len(e) {
		return e[len(e):]
	}
	return e[n:]
}

func NextBlock(e Expression) (Expression, int) {
	return NextSegment(e, "{", "}")
}

func ExtractParas(e Expression) (Expression, int) {
	return NextSegment(e, "(", ")")
}

func NextSegment(e Expression, start, end string) (Expression, int) {
	if e[0] != start {
		panic(fmt.Sprintf("Invalid block: %v", e))
	}
	var count = 1
	for i := 1; i < len(e); i++ {
		if e[i] == end {
			count--
			if count == 0 {
				return e[1:i], i
			}
		}
		if e[i] == start {
			count++
		}
	}
	panic(fmt.Sprintf("unenclosed block: %v", e))
}

type Procedure func(...interface{}) interface{}

func IsPrimitive(app interface{}) bool {
	_, ok := app.(Procedure)
	return ok
}

func Primitive(app interface{}) Procedure {
	return app.(Procedure)
}

type CompoundProcedure struct {
	paras []string
	body  Expression
	env   *Environment
}

func NewCompoundProdedure(paras []string, body Expression, env *Environment) CompoundProcedure {
	return CompoundProcedure{
		paras: paras,
		body:  body,
		env:   env,
	}
}

func IsCompound(app interface{}) bool {
	_, ok := app.(CompoundProcedure)
	return ok
}

func Compound(app interface{}) CompoundProcedure {
	return app.(CompoundProcedure)
}

func (c *CompoundProcedure) Body() Expression {
	return c.body
}

func (c *CompoundProcedure) Paras() []string {
	return c.paras
}

func (c *CompoundProcedure) Env() *Environment {
	return c.env
}
