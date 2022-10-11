package main

import (
	"context"
	"fmt"
	"time"
)

type SlowFunction func(string) (string, error)

type WithContext func(context.Context, string) (string, error)

func demo(s string) (string, error) {
	time.Sleep(time.Second * 10)
	return s, nil
}

func Timeout(sf SlowFunction) WithContext {
	return func(ctx context.Context, s string) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		chres := make(chan string)
		cherr := make(chan error)
		go func() {
			res, err := sf(s)
			chres <- res
			cherr <- err
		}()

		select {
		case res := <-chres:
			return res, <-cherr
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
func main() {
	timeout := Timeout(demo)

	ctxShortTimeout, cancelShortTimeout := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelShortTimeout()
	r, err := timeout(ctxShortTimeout, "my test")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

	ctxLongTimeout, cancelLongTimeout := context.WithTimeout(context.Background(), time.Second*20)
	defer cancelLongTimeout()
	r, err = timeout(ctxLongTimeout, "my test")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}
