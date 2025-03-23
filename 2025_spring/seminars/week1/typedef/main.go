package main

import "fmt"

type A struct {
	Value int
}

func (a *A) F() string {
	return "a"
}

type B A

type UserID int

// type UserID = int

func G(id UserID) {
	fmt.Println(id)

	asdklfsa

	defer f()

	jksladfjldska

	if asdf {
		return
	}

	defer g()
}



func main() {
	var v int = 5
	t := int(v)
	G(UserID(v))
}

func main() {
	c := B{Value: 5}
	fmt.Println(c.Value)
	t := A(c)
	fmt.Println(t.F())
}
