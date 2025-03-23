package main

import "fmt"

func main() {
	a := 0

	for i := 0; i < 2; i++ {
		defer func(b *int) {
			fmt.Println("here ", i, *b)
		}(&a)
	}

	a += 1

	if a == 1 {
		return
	}

	defer func() {
		fmt.Println("here3")
	}()

}
