package main

import "fmt"

type A struct{}

// func (a *A) f() string {
// 	return "a"
// }

// type B struct {
// 	A
// }

// type C struct {
// 	A
// }

// type D struct {
// 	B
// 	C
// }

type B struct {
	A
}

func (b *B) f() string {
	return "b"
}

type C struct {
	A
}

func (c *C) f() string {
	return "c"
}

type D struct {
	B
	C
}

func (d *D) f() string {
	return d.C.f()
}

func main() {
	d := D{}
	fmt.Println(d.f())
}
