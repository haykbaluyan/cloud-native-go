package main

import (
	"fmt"
	"sync"
)

func Split(source <-chan int, n int) []<-chan int {
	dests := make([]<-chan int, 0)

	for i := 0; i < n; i++ {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for val := range source {
				ch <- val
			}
		}()
		dests = append(dests, ch)
	}

	return dests
}
func main() {
	source := make(chan int)
	dests := Split(source, 5)
	go func() {
		for i := 1; i <= 10; i++ {
			source <- i
		}
		close(source)
	}()

	var wg sync.WaitGroup
	wg.Add(len(dests))

	for i, ch := range dests {
		go func(i int, c <-chan int) {
			for val := range c {
				fmt.Printf("val %d from channel %d\n", val, i)
			}
			wg.Done()
		}(i, ch)
	}

	wg.Wait()
}
