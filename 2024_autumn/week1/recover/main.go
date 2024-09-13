package main

import "fmt"

func g() {
	fmt.Println("here1")
	panic("PRESS F")
	fmt.Println("here2")
}

func f() {
	a := 5
	defer func() {
		fmt.Print("aaa1\n")
		if r := recover(); r != nil {
			fmt.Println("Recovered: ", r)
		}
	}()
	defer func() {
		fmt.Printf("aaa2: %d\n", a)
	}()

	a = 6
	fmt.Println("here3")
	g()
	fmt.Println("here4")
}

func main() {
	f()
}
