package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	if len(channels) == 0 {
		return nil
	}
	done := make(chan interface{})
	var once sync.Once

	for _, c := range channels {
		go func(c <-chan interface{}) {
			select {
			case <-c:
				once.Do(func() { close(done) })
			case <-done:

			}
		}(c)
	}
	return done
}

func orSecond(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	done := make(chan interface{})
	go func() {
		defer close(done)
		select {
		case <-channels[0]:
		case <-channels[1]:
		case <-orSecond(append(channels[2:], done)...):
		}
	}()
	return done
}
