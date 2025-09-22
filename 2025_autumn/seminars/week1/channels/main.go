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
		num := rand.Intn(100)
		fmt.Println("Sending ", num)
		ch <- num
		fmt.Println("Sent ", num)
	}
}

func logNumbers(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		time.Sleep(1 * time.Second)
		fmt.Println("Received number ", num)
	}
}

func main() {
	ch := make(chan int, 3)

	wg1 := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg1.Add(1)
		go generateNumbers(ch, &wg1)
	}

	for i := 0; i < 3; i++ {
		wg2.Add(1)
		go logNumbers(ch, &wg2)
	}

	wg1.Wait()
	close(ch)
	wg2.Wait()
	fmt.Println("HERE")
}
