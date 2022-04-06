[![Latest Release](https://img.shields.io/github/release/nichtsen/goeval.svg)](https://github.com/nichtsen/goeval/releases)
[![Go ReportCard](https://goreportcard.com/badge/nichtsen/goeval)](https://goreportcard.com/report/nichtsen/goeval)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/nichtsen/goeval)

A tiny evaluator of scheme expressions 

![demo](https://github.com/nichtsen/lis/raw/main/eval/img/demo.gif)

## expression
- define-expr
```
define "{" 
    var value-expr | 
    proc"(" [arg,] ")" 
    "{" 
         expr [, expr]
    "}" 
"}"
```
- assign-expr
```
set var value-expr
```
- if-expr
```
if predict-expr
    "{" 
    expr [, expr]
    "}"

```
- lambda-expr
```
lambda"("[arg,] ")" "{" expr [,expr] "}" ["(" [,para] ")"]
```

- application-expr
```
proc-name"("[,para]")" 
```
- perform-expr
```
perform application-expr
```

## Remark
Core of evaluator is interleaving of `eval` and `apply` loop. Everytime eval enconters a
application it call apply to make actually calling.
A compound procedure is a sequence of expression to be evaluated 
within a given environment.

## Example 
compute fibonacci with evaluator:
```eval
Eval-go INPUT:
define fib(n) { 
    if ==(n, 0) {0} 
    if ==(n, 1) {1} 
    +(fib(-(n,1)), fib(-(n,2))) 
    };
Eval-go INPUT:
fib(10);
Eval-go OUTPUT:
 55 
```
# document
[godoc](https://godoc.org/github.com/nichtsen/goeval)
