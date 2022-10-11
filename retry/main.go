package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type retryWrapper func(context.Context, string) (string, error)

type effector func(string) (string, error)

func demo() func(string) (string, error) {
	count := 0
	return func(s string) (string, error) {
		count++
		if count < 5 {
			return "", errors.New("something went wrong")
		}
		return s, nil
	}
}

func retry(e effector, retryLimit int, delay time.Duration) retryWrapper {
	wrapper := func(ctx context.Context, s string) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		var err error
		for i := 0; i < retryLimit; i++ {
			res, err := e(s)
			if err == nil {
				return res, nil
			}

			fmt.Printf("Attempt %d failed; retrying in %v\n", i+1, delay)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}

		return "", err
	}
	return wrapper

}
func main() {
	f := retry(demo(), 10, time.Second*2)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		ticker := time.NewTicker(time.Second * 30)
		<-ticker.C
		ticker.Stop()
	}()

	r, err := f(ctx, "my test")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}
