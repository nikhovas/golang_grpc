package main

import "fmt"

type A struct{}

func (a *A) F() string {
	return "a"
}

type B struct{}

// func (b *B) F() string {
// 	return "b"
// }

type C struct {
	A
	B
}

// func (a *C) F() string {
// 	return "c"
// }

func main() {
	c := C{}
	fmt.Println(c.F())
}
