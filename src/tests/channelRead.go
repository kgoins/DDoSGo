package main

import (
	"fmt"
	"time"
)

type Worker struct {
	quit chan bool
}

func newWorker(quit chan bool) Worker {
	return Worker{quit: quit}
}

func (worker Worker) Run() {
	go func() {
		fmt.Println("worker starting")
		for {
			select {
			case <-worker.quit:
				fmt.Println("worker closing")
				return
			}
		}
	}()
}

func main() {
	numWorkers := 10
	quit := make(chan bool)

	for i := 0; i <= numWorkers; i++ {
		worker := newWorker(quit)
		worker.Run()
		fmt.Printf("worker %d created\n", i)
	}

	time.Sleep(1000 * time.Millisecond)

	fmt.Println("sending killsig")
	quit <- true

	time.Sleep(1000 * time.Millisecond)
	fmt.Println("all workers dead")
}
