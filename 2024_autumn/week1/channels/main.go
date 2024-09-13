package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func generateNumber(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		num := rand.Intn(100)
		fmt.Printf("Start generating: %d\n", num)
		ch <- num
		fmt.Printf("Generated: %d\n", num)
		// time.Sleep(time.Millisecond * 100)
	}
}

func logNumbers(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Printf("Logged: %d\n", num)
		time.Sleep(time.Millisecond * 300)
	}
}

func main() {
	ch := make(chan int, 100)
	var wg sync.WaitGroup

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go generateNumber(ch, &wg)
	}

	var wg2 sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg2.Add(1)
		go logNumbers(ch, &wg2)
	}

	wg.Wait()
	close(ch)
	wg2.Wait()
}
