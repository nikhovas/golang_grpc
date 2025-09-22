package main

import "fmt"

func f() {
	var k int = 2

	defer func(y *int) {
		fmt.Println(*y)
	}(&k)

	k = 4
}

func main() {
	f()
}
