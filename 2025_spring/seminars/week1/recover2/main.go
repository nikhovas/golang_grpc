package main

import (
	"fmt"
	"sync"
)

func g(wg *sync.WaitGroup) {
	defer wg.Done()
	// defer func() {
	// 	fmt.Println("aaa")
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Recovered: ", r)
	// 	}
	// }()
	y := 0
	y += 5
	var g [3]int
	fmt.Println(g[y])
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go g(&wg)
	// do smth .....
	wg.Wait()
	fmt.Println("here1")
}
