package main

import "fmt"
import "os"
import "github.com/amc90/topsy"

func main() {
	fmt.Println("Hello, world!")
	lex, e:=topsy.Lex(os.Stdin)
	fmt.Println(e)
	fmt.Println(lex)
}

