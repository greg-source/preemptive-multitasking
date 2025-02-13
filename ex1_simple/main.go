package main

import (
	"fmt"
	"runtime"
	"time"
)

func tightLoop() {
	fmt.Println("Goroutine 1 starting")
	for {
		// Tight loop with no function calls or yielding
	}
}

func main() {
	// Limit the scheduler to 1 OS thread
	runtime.GOMAXPROCS(1)

	// Start a long-running goroutine with a tight loop
	go tightLoop()

	// Start another goroutine
	go func() {
		for {
			fmt.Println("Goroutine 2 running")
			time.Sleep(1 * time.Second)
		}
	}()

	// Give the goroutines some time to run
	time.Sleep(3 * time.Second)
	fmt.Println("Main function done")
}
