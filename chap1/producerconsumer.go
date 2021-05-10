package chap1

import (
	"fmt"
	"time"
)

func product(base int, out chan<- int) {
	for i := 0; i < 100; i++ {
		out <- i * base
	}
}

func consume(in <-chan int) {
	for value := range in {
		fmt.Printf("%d ", value)
	}
	fmt.Println()
}

func ProductConsume() {
	queue := make(chan int)
	go product(3, queue)
	go product(5, queue)
	go product(7, queue)
	go consume(queue)
	time.Sleep(2 * time.Second)
}
