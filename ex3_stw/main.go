package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

const (
	NumTasks      = 500   // Number of tasks.
	MemoryIntense = 10000 // Size of memory-intensive task (number of elements).
)

func complexOperation(wg *sync.WaitGroup) int64 {
	defer wg.Done()

	var result int64 = 1
	const iterations int64 = 500_000_000 // Adjust this value as needed
	for i := int64(0); i < iterations; i++ {
		// Pure arithmetic only—no function calls!
		result = result*3 + 7 - result%5
	}
	return result
}

func main() {
	// Create trace file and handle any errors.
	f, err := os.Create("trace.out")
	if err != nil {
		fmt.Println("Error creating trace file:", err)
		return
	}
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()

	// Task queue and result queue.
	taskQueue := make(chan int, NumTasks)
	resultQueue := make(chan int, NumTasks)

	// Launch 4 workers to perform Memory Intensive Tasks
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(taskQueue, resultQueue)
		}()
	}

	// Launch complex operation (which takes ~2 seconds)
	wg.Add(1)
	go complexOperation(&wg)

	// Send tasks to the taskQueue and close after sending
	waitToFinish(taskQueue, resultQueue, &wg)

	// Process the results.
	for result := range resultQueue {
		_ = result
	}

	// All work is finished.
	fmt.Println("Done!")
}

func waitToFinish(taskQueue, resultQueue chan int, wg *sync.WaitGroup) {
	// Send tasks to the queue.
	for i := 0; i < NumTasks; i++ {
		taskQueue <- i
	}
	close(taskQueue)

	// Wait for all workers to finish before closing the resultQueue.
	wg.Wait()
	close(resultQueue)
}

// Worker function.
func worker(tasks <-chan int, results chan<- int) {
	for task := range tasks {
		result := performMemoryIntensiveTask(task)
		results <- result
	}
}

// performMemoryIntensiveTask is a memory-intensive function.
func performMemoryIntensiveTask(task int) int {
	// Create a large-sized slice.
	data := make([]int, MemoryIntense)
	for i := 0; i < MemoryIntense; i++ {
		data[i] = i + task
	}

	// Latency imitation.
	time.Sleep(10 * time.Millisecond)

	// Calculate the result by summing the slice.
	result := 0
	for _, value := range data {
		result += value
	}
	return result
}
