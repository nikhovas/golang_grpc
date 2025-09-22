package main

import "fmt"

type A struct {
	Value int
}

func (a *A) f() string {
	return "a"
}

type B A

type C = A

// func main() {
// 	b := B{}
// 	y := b.Value
// 	fmt.Println(y)
// 	fmt.Println(b.f())
// }

func main() {
	c := C{}
	y := c.Value
	fmt.Println(y)
	fmt.Println(c.f())
}
