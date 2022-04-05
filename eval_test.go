package eval

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleEnvironment() {
	vars := []string{"a", "b", "c"}
	val := []interface{}{"a", 3.14, 1}

	env := ExtendEnv(GlobalEnv, vars, val)
	env.SetVariable("b", "b")
	env.DefineVariable("d", "d")
	vars = append(vars, "d")
	for _, va := range vars {
		val := env.LookUpVariable(va)
		fmt.Println(val)
	}
	// Output:
	// a
	// b
	// 1
	// d
}

func TestProdedureExpr(t *testing.T) {
	InitGlobal()
	expr := Expression([]string{"foo", "(", "a,b,c", ")", "{", "set", "a", "1", "a+b+c", "}"})
	if !ProcedureExpr(expr) {
		t.Error("failed prediction of procedure expression")
	}
	va := ProcedureVar(expr)
	if va != "foo" {
		t.Errorf("expected \"foo\", got %v", va)
	}
	paras, idx := ProcedureParas(expr)
	ep := []string{"a", "b", "c"}
	if !reflect.DeepEqual(paras, ep) {
		t.Errorf("expected to be equal, %v vs %v", paras, ep)
	}

	eb := expr[5:9]
	body, _ := ProcedureBody(expr[idx:])

	if !reflect.DeepEqual(body, eb) {
		t.Errorf("expected to be equal, %v vs %v", body, eb)
	}
}

func TestBraces(t *testing.T) {
	InitGlobal()
	text := `define foo(a){{+(a,1)} }`
	expr := MakeExpr(text)
	val := fmt.Sprintf("%s", *expr)
	expected := "[define foo ( a ) { { + ( a , 1 ) } }]"
	if val != expected {
		t.Errorf("expected to be %v, not %v", expected, val)
	}
}

func TestEval01(t *testing.T) {
	InitGlobal()
	text := `define b 1 b`
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != 1 {
		t.Errorf("expected to be 1, not %v", val)
	}

	text = `set b 2 b`
	expr = MakeExpr(text)
	val = Eval(expr, GlobalEnv)
	if val != 2 {
		t.Errorf("expected to be 2, not %v", val)
	}
}

func TestEval02(t *testing.T) {
	InitGlobal()
	text := `define a 1 define foo(v) { v } foo(a) `
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != 1 {
		t.Errorf("expected to be 1, not %v", val)
	}
}

func TestEvalIf01(t *testing.T) {
	InitGlobal()
	text := `define a 1 define b 2 if >(b,a) { b } a `
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != 2 {
		t.Errorf("expected to be 2, not %v", val)
	}
}

func TestEvalIf02(t *testing.T) {
	InitGlobal()
	text := `define a 2 define b 1 if >(b,a) { b } a `
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != 2 {
		t.Errorf("expected to be 2, not %v", val)
	}
}

func TestEvalSymbol(t *testing.T) {
	InitGlobal()
	text := `define a 'hello a`
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != "hello" {
		t.Errorf("expected to be hello, not %v", val)
	}
}

func TestLambda01(t *testing.T) {
	InitGlobal()
	text := `define square 'undefined set square lambda(a) { *(a,a) } square(2)  `
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != 4 {
		t.Errorf("expected to be 4, not %v", val)
	}
}

func TestLambda02(t *testing.T) {
	InitGlobal()
	// double is a procedure that return a procedure
	text := `define double(s) { lambda(prefix) { append(prefix,append(s,s)) } } define proc double('z)  proc('doublez:)`
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != "doublez:zz" {
		t.Errorf("expected to be \"doublez:zz\", not %v", val)
	}
}

func TestLambda03(t *testing.T) {
	InitGlobal()
	text := `lambda(a){*(a,a)}(4)`
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != 16 {
		t.Errorf("expected to be 16, not %v", val)
	}
}

func TestLambda04(t *testing.T) {
	InitGlobal()
	text := `define square(a){*(a,a)} square(lambda(a){*(a,a)}(2))`
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != 16 {
		t.Errorf("expected to be 16, not %v", val)
	}
}
func TestFactorial01(t *testing.T) {
	InitGlobal()
	text := `define fact(n) { if ==(n,1) { 1 } define tmp -(n,1) *(n,fact(tmp)) } fact(3)`
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != 6 {
		t.Errorf("expected to be 6, not %v", val)
	}
}

func TestFactorial02(t *testing.T) {
	InitGlobal()
	text := `define fact(n) { if ==(n,1) { 1 } *(n,fact(-(n,1))) } fact(3)`
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != 6 {
		t.Errorf("expected to be 6, not %v", val)
	}
}

func TestFibonacci(t *testing.T) {
	InitGlobal()
	text := `
	   define fib(n){ 
		   if ==(n, 0) {0} 
		   if ==(n, 1) {1} 
		   +(fib(-(n,1)),fib(-(n,2))) 
		   };
		fib(6) `
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != 8 {
		t.Errorf("expected to be 8, not %v", val)
	}
}

func TestCons01(t *testing.T) {
	InitGlobal()
	text := `
	   define a cons('a,'b)
	   car(a) `
	expr := MakeExpr(text)
	val := Eval(expr, GlobalEnv)
	if val != "a" {
		t.Errorf("expected to be a, not %v", val)
	}
}

func ExamplePerformExpr() {
	InitGlobal()
	text := `
	   define l list('a,'b,'c)
	   define printlist(l) {
		   perform print(car(l))
		   if not-null?(cdr(l)) {
				printlist(cdr(l))
		   }
	   }
	 	printlist(l)  
	   `
	expr := MakeExpr(text)
	Eval(expr, GlobalEnv)
	//output:
	// abc
}

func ExampleList() {
	InitGlobal()
	text := `
	   define tab list(list(1,'a),list(2,'b),list(3,'c))
	   define printlist(l) {
		   define cur car(l)
		   perform print(car(cur))
		   perform print(car(cdr(cur)))
		   if not-null?(cdr(l)) {
				printlist(cdr(l))
		   }
	   }
	 	printlist(tab)
	   `
	expr := MakeExpr(text)
	Eval(expr, GlobalEnv)
	// Output:
	// 1a2b3c

}
