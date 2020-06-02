package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var x int64
var wg sync.WaitGroup

func calc() {
	for i := 0; i < 5000; i++ {
		atomic.AddInt64(&x, 1)
	}

	wg.Done()
}

func main() {

	wg.Add(2)
	go calc()
	go calc()
	wg.Wait()

	fmt.Println(x)
}
