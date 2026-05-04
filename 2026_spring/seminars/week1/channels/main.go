package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func generateNumbers(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		num := rand.Intn(100)
		fmt.Println("Sending ", num)
		ch <- num
		fmt.Println("Sent ", num)
	}
}

func logNumbers(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Println("Received number ", num)
	}
}

func main() {
	ch := make(chan int, 10000)

	wg1 := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}

	wg1.Add(5)
	for i := 0; i < 5; i++ {
		go generateNumbers(ch, &wg1)
	}

	wg2.Add(3)
	for i := 0; i < 3; i++ {
		go logNumbers(ch, &wg2)
	}

	wg1.Wait()
	close(ch)
	wg2.Wait()
	fmt.Println("HERE")
}

// func main() {
// 	ch := make(chan int, 1)
// 	ch <- 5
// 	close(ch)
// 	ch <- 7
// 	n := <-ch
// 	y, b := <-ch
// 	fmt.Println(n)
// 	fmt.Println(y)
// 	fmt.Println(b)
// }
