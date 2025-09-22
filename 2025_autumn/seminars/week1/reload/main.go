package main

import "fmt"

type A struct{}

func (a *A) f() string {
	return "a"
}

type B struct {
	A
}

type C struct {
	A
}

type D struct {
	B
	C
}

// func (b *B) F() string {
// 	return "b"
// }

// type C struct {
// 	A
// 	B
// }

func main() {
	d := D{}
	fmt.Println(d.f())
}
