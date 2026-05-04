package main

import "fmt"

type A struct {
	Value int
}

func (a *A) f() string {
	return "a"
}

// type B A

type B = A

// func main() {
// 	b := B{}

// 	a := A(b)
// 	fmt.Println(a.f())
// 	// fmt.Println(ggg.Value)
// 	// fmt.Println(A(b).f())
// }

func main() {
	b := B{}

	fmt.Println(b.f())
	// fmt.Println(ggg.Value)
	// fmt.Println(A(b).f())
}
