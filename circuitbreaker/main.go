package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Breaker func(context.Context, string) (string, error)

type Circuit func(string) (string, error)

func demo() func(string) (string, error) {
	count := 0
	return func(s string) (string, error) {
		count++
		if count < 4 {
			return "", errors.New("something went wrong")
		}

		return s, nil
	}
}

func CircuitBreaker(circuit Circuit, maxFailures int, d time.Duration) Breaker {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context, s string) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		m.RLock()
		if consecutiveFailures >= maxFailures {
			shouldRetryAt := lastAttempt.Add(d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		m.RUnlock()
		res, err := circuit(s)

		m.Lock()
		defer m.Unlock()
		lastAttempt = time.Now()
		if err != nil {
			consecutiveFailures++
			return "", err
		}

		consecutiveFailures = 0
		return res, nil
	}
}

func main() {
	f := CircuitBreaker(demo(), 2, time.Second*5)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		ticker := time.NewTicker(time.Second * 45)
		<-ticker.C
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
