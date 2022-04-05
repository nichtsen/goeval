package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nichtsen/lis/eval"
)

func main() {
	s := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Eval-go INPUT:")
		r, err := s.ReadString(byte(';'))
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
		str := strings.ToLower(string(r))
		if str == "q" || str == "quit" || str == "exit" {
			break
		}
		expr := eval.MakeExpr(string(r))
		val := eval.Eval(expr, eval.GlobalEnv)
		if val == nil {
			continue
		}
		fmt.Printf("Eval-go OUTPUT:\n %v \n", val)
	}
}
