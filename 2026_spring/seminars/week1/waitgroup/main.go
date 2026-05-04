package main

import (
	"fmt"
	"sync"
)

func ExternalSystemA() int {
	return 1
}

func ExternalSystemB() int {
	return 2
}

// func GrpcHandler() {
// 	wg := sync.WaitGroup{}
// 	wg.Add(2)
// 	var aRes, bRes int
// 	go func() {
// 		defer wg.Done()
// 		aRes = ExternalSystemA()
// 	}()
// 	go func() {
// 		defer wg.Done()
// 		bRes = ExternalSystemB()
// 	}()
// 	wg.Wait()
// 	_, _ = fmt.Printf("%d %d\n", aRes, bRes)
// }

func GrpcHandler() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var aRes, bRes int
	go func() {
		defer wg.Done()
		aRes = ExternalSystemA()
	}()
	bRes = ExternalSystemB()
	wg.Wait()
	_, _ = fmt.Printf("%d %d\n", aRes, bRes)
}

func main() {
	GrpcHandler()
}
