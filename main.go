package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	goroutines = 5
	jobs       = 20
)

func main() {
	normal()
	outro()
}

func work(id int, query string) string {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("job %d = %s\n", id, query)
}

func normal() {
	start := time.Now()
	inExecution := int32(0)

	var results []string
	var wg sync.WaitGroup
	ticket := make(chan struct{}, goroutines)
	for i := 0; i < jobs; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ticket <- struct{}{}
			atomic.AddInt32(&inExecution, 1)
			s := work(i, fmt.Sprintf("string%d - goroutines: %d", i, inExecution))
			results = append(results, s)
			<-ticket
			atomic.AddInt32(&inExecution, -1)
		}(i)
	}
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}

func outro() {
	start := time.Now()

	var results []string
	var ticket = make(chan bool, goroutines)
	c := make(chan string)
	for i := 0; i < jobs; i++ {
		go func(i int) {
			ticket <- true
			defer func() { <-ticket }()
			c <- work(i, fmt.Sprintf("string%d - goroutines: %d", i, len(ticket)))
		}(i)
	}

	for i := 0; i < jobs; i++ {
		results = append(results, <-c)
	}

	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
