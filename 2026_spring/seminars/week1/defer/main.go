package main

import "fmt"

func f() {
	defer func() {
		fmt.Println("A")
	}()

	defer func() {
		fmt.Println("C")
	}()

	fmt.Println("B")

	return
}

func main() {
	f()
}
