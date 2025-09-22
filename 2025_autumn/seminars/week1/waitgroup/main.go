package main

import (
	"fmt"
	"sync"
	"time"
)

func externalSystemCall1(wg *sync.WaitGroup) int {
	defer wg.Done()
	time.Sleep(1 * time.Second)
	fmt.Println("HERE1")
	return 1
}

func externalSystemCall2(wg *sync.WaitGroup) string {
	defer wg.Done()
	time.Sleep(2 * time.Second)
	fmt.Println("HERE2")
	return "aaa"
}

func externalSystemCall3(wg *sync.WaitGroup) bool {
	// defer wg.Done()
	time.Sleep(3 * time.Second)
	fmt.Println("HERE3")
	return true
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var intVal int
	var stringVal string
	var boolVal bool
	go func() {
		intVal = externalSystemCall1(&wg)
	}()
	go func() {
		stringVal = externalSystemCall2(&wg)
	}()

	boolVal = externalSystemCall3(&wg)

	wg.Wait()

	fmt.Println(intVal, stringVal, boolVal)
}
