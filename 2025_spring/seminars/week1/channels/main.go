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
		fmt.Println("Received number", num)
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(1)
	for i := 0; i < 1; i++ {
		go generateNumbers(ch, &wg)
	}

	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	for i := 0; i < 1; i++ {
		go logNumbers(ch, &wg2)
	}

	wg.Wait()
	close(ch)
	wg2.Wait()
}
