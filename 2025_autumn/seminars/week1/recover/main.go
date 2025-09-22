package main

import "fmt"

func g() {
	fmt.Println("here1")
	panic("F")
	fmt.Println("here2")
}

func f() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered ", r)
		}
	}()

	g()
}

func main() {
	f()
}
