package dispatcher

import "fmt"
import "sync"

type Worker struct {
	workerPool         chan chan Dispatchable
	inboundWorkChannel chan Dispatchable
	quit               chan bool
	waitgroup          *sync.WaitGroup
}

func NewWorker(workerPool chan chan Dispatchable) *Worker {
	inboundWorkChannel := make(chan Dispatchable)
	quit := make(chan bool)

	return &Worker{
		workerPool:         workerPool,
		inboundWorkChannel: inboundWorkChannel,
		waitgroup:          &sync.WaitGroup{},
		quit:               quit}
}

func (worker *Worker) Close() {
	worker.quit <- true
	worker.waitgroup.Wait()
}

func (worker *Worker) Start() {
	worker.waitgroup.Add(1)
	go func() {
		defer worker.waitgroup.Done()
		for {
			select {
			case <-worker.quit:
				return

			case worker.workerPool <- worker.inboundWorkChannel:
				work, closed := <-worker.inboundWorkChannel
				if !closed {
					fmt.Println("worker pool was closed")
					return
				}

				work.DispatchableExec()
			}
		}
	}()
}
