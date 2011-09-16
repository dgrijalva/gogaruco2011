package main

import (
	"fmt"
	"time"
	"rand"
)

var done = make(chan int, 10)

func main() {
	c := make(chan int)
	go produce(c, 1000)
	// start consumers
	for i := 0; i < 10; i++ {
		go consume(c, i)
	}
	// wait for consumers to finish
	for i := 0; i < 10; i++ {
		<-done
	}
}

func produce(c chan int, until int) {
	var a,b int
	c <- 0
	// c <- 1
	for a, b = 0, 1; b < until; a, b = b, a+b {
		c <- b
	}
	close(c)
}

func consume(c chan int, n int) {
	for {
		time.Sleep(rand.Int63n(100) * 1e6)
		if i, ok := <-c; ok {
			fmt.Printf("Consumer %v consumed %v\n", n, i)
		} else {
			done <-1
			return
		}
	}
}