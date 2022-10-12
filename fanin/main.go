package main

import (
	"fmt"
	"sync"
	"time"
)

func Funnel(sources ...<-chan int) <-chan int {
	dest := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(sources))

	for _, ch := range sources {
		go func(c <-chan int) {
			defer wg.Done()
			for n := range c {
				dest <- n
			}
			fmt.Println("channel closed")
		}(ch)
	}

	go func() {
		wg.Wait()
		close(dest)
	}()

	return dest
}
func main() {
	sources := make([]<-chan int, 0)
	for i := 0; i < 3; i++ {
		ch := make(chan int)
		sources = append(sources, ch)
		go func() {
			defer close(ch)
			for i := 0; i <= 2; i++ {
				ch <- i
				time.Sleep(time.Second * 1)
			}
		}()
	}

	dest := Funnel(sources...)
	for d := range dest {
		fmt.Println(d)
	}
}
