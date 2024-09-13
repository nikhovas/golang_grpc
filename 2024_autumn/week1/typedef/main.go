package main

import "fmt"

type Number interface {
	Square() int
}

func f(a Number) {
	fmt.Printf("%d", a.Square())
}

type A int

func (a A) Square() int {
	return int(a) * int(a)
}

func main() {
	a := 5
	f(A(a))
}
