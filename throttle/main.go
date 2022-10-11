package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type effector func(context.Context, string) (string, error)

type effectee func(string) (string, error)

func demo(s string) (string, error) {
	return s, nil
}

func throttle(effee effectee, max int, d time.Duration) effector {
	count := 0
	var once sync.Once
	effer := func(ctx context.Context, s string) (string, error) {
		if ctx.Err() != nil {
			fmt.Println("context here")
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

		res, err := effee(s)
		if err != nil {
			return "", err
		}

		return res, nil
	}
	return effer

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
