package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type throttleWrapper func(context.Context, string) (string, error)

type effector func(string) (string, error)

func demo(s string) (string, error) {
	return s, nil
}

func throttle(e effector, max int, d time.Duration) throttleWrapper {
	count := 0
	var once sync.Once
	return func(ctx context.Context, s string) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		once.Do(func() {
			ticker := time.NewTicker(d)
			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						count = 0
					}
				}
			}()
		})

		if count >= max {
			return "", errors.New("too many calls, waiting for some time to pass")
		}

		count++

		res, err := e(s)
		if err != nil {
			return "", err
		}

		return res, nil
	}
}

func main() {
	f := throttle(demo, 2, time.Second*10)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		ticker := time.NewTicker(time.Second * 30)
		<-ticker.C
		ticker.Stop()
	}()

	for {
		if ctx.Err() != nil {
			return
		}
		r, err := f(ctx, "my test")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(r)
		}

		time.Sleep(time.Second * 3)
	}

}
