# eval
A tiny evaluator of scheme expressions  

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
application it call apply to make acctual calling.
A compound procedure is a sequence of expression to be evaluated 
within a given environment.

## Example 
compute fibonacci with evaluator:
```eval
Eval-go INPUT:
define fib(n) { 
    if ==(n,0) { 0 } 
    if ==(n,1) { 1 } 
    +(fib(-(n,1)),fib(-(n,2))) 
    };
Eval-go INPUT:
fib(10);
Eval-go OUTPUT:
 55 
```
