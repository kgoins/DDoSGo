package dispatcher

import "fmt"

type Dispatchable interface {
	DispatchableExec()
}

type Dispatcher struct {
	maxWorkers    int
	inboundWork   chan Dispatchable
	workerChannel chan chan Dispatchable
	workerPool    []*Worker
	shutdown      chan bool
}

func NewDispatcher(inboundWork chan Dispatchable, workers int) *Dispatcher {
	workerChannel := make(chan chan Dispatchable)
	shutdown := make(chan bool)

	return &Dispatcher{
		maxWorkers:    workers,
		inboundWork:   inboundWork,
		workerChannel: workerChannel,
		workerPool:    make([]*Worker, workers),
		shutdown:      shutdown}
}

func (dispatcher *Dispatcher) Close() {
	// kill dispatcher goroutine
	dispatcher.shutdown <- true

	for _, worker := range dispatcher.workerPool {
		worker.Close()
	}

	close(dispatcher.workerChannel)
	close(dispatcher.shutdown)
	fmt.Println("all workers shutdown")
}

func (dispatcher *Dispatcher) Run() {

	for i := 0; i < dispatcher.maxWorkers; i++ {
		worker := NewWorker(dispatcher.workerChannel)
		dispatcher.workerPool[i] = worker
		worker.Start()
	}

	go dispatcher.dispatch()
}

func (dispatcher *Dispatcher) dispatch() {
	for {
		select {
		case <-dispatcher.shutdown:
			fmt.Println("dispatch goroutine closing")
			return

		case work := <-dispatcher.inboundWork:
			availWorker := <-dispatcher.workerChannel
			availWorker <- work
			fmt.Println("received work")
		}
	}
}
