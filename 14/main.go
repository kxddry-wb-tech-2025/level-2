package main

import (
	"fmt"
	"sync"
	"time"
)

func or(channels ...<-chan any) <-chan any {
	out := make(chan any)
	var once sync.Once
	stop := make(chan struct{})

	closeDone := func() {
		once.Do(func() { close(stop) })
	}

	for _, c := range channels {
		go func(c <-chan any) {
			for {
				select {
				case v, ok := <-c:
					if !ok {
						closeDone()
						return
					}
					select {
					case out <- v:
					case <-stop:
						return
					}
				}
			}
		}(c)
	}

	go func() {
		<-stop
		close(out)
	}()

	return out
}

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
	fmt.Printf("done after %v", time.Since(start)) // done after 1.001203083s
}
