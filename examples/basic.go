package main

import (
	"fmt"
	"github.com/opesun/jsonlang"
)

var src = `
isOdd(&is_odd, 3);
set(&x, { "val": is_odd });
set(&f, {"lol":{"trolol":22}});
println(f.lol.trolol);
println(x);
`

func isOdd(a ...interface{}) {
	if len(a) < 2 {
		panic("Too few parameters at isOdd.")
	}
	num := int(a[1].(float64))
	verdict := num%1 == 0
	a[0].(jsonlang.Ref).Set(verdict)
}

type m map[string]func(...interface{})

func main() {
	src, err := jsonlang.Compile(src)
	if err != nil {
		panic(err)
	}
	err = jsonlang.Interpret(src, nil, m{"isOdd": isOdd})
	fmt.Println("Error:", err)
}
