package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func SlowOperation(s string) string {
	time.Sleep(time.Second * 5)
	return s
}

func ConcurentSlowOperation(ctx context.Context, s string) <-chan string {
	out := make(chan string)
	go func() {
		res := SlowOperation(s)
		out <- res
		close(out)
	}()

	return out
}

type Future interface {
	Result() (string, error)
}

type InnerFuture struct {
	once sync.Once
	wg   sync.WaitGroup

	resCh <-chan string
	errCh <-chan error

	res string
	err error
}

func (f *InnerFuture) Result() (string, error) {
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()
		f.res = <-f.resCh
		f.err = <-f.errCh
	})

	f.wg.Wait()

	return f.res, f.err
}

func SlowOperationWithFuture(ctx context.Context, s string) *InnerFuture {
	resCh := make(chan string)
	errCh := make(chan error)

	go func() {
		select {
		case <-time.After(time.Second * 5):
			resCh <- s
		case <-ctx.Done():
			errCh <- ctx.Err()
		}
		close(resCh)
		close(errCh)
	}()

	return &InnerFuture{
		resCh: resCh,
		errCh: errCh,
	}
}

func main() {
	/*c1 := ConcurentSlowOperation("test1")
	c2 := ConcurentSlowOperation("test2")
	c3 := ConcurentSlowOperation("test3")
	fmt.Println("started multiple operations")
	fmt.Printf("%s %s %s\n", <-c1, <-c2, <-c3)*/

	ctx := context.Background()
	f1 := SlowOperationWithFuture(ctx, "test1")
	f2 := SlowOperationWithFuture(ctx, "test2")
	f3 := SlowOperationWithFuture(ctx, "test3")
	fmt.Printf(f1.Result())
	fmt.Printf(f2.Result())
	fmt.Printf(f3.Result())

}
