package dispatcher

import "net"
import "fmt"

import "msgs"

type Dispatcher struct {
	maxWorkers    int
	inboundConns  chan net.Conn
	workerChannel chan chan net.Conn
	workerPool    []*MsgWorker
	shutdown      chan bool
}

func NewDispatcher(inboundConns chan net.Conn, workers int) *Dispatcher {
	workerChannel := make(chan chan net.Conn)
	shutdown := make(chan bool)

	return &Dispatcher{
		maxWorkers:    workers,
		inboundConns:  inboundConns,
		workerChannel: workerChannel,
		workerPool:    make([]*MsgWorker, workers),
		shutdown:      shutdown}
}

func (dispatcher *Dispatcher) Close() {
	// kill dispatcher goroutine
	dispatcher.shutdown <- true

	for _, worker := range dispatcher.workerPool {
		worker.Close()
	}

	close(dispatcher.workerChannel)
	fmt.Println("all workers shutdown")
}

func (dispatcher *Dispatcher) Run() {

	for i := 0; i < dispatcher.maxWorkers; i++ {
		worker := NewMsgWorker(dispatcher.workerChannel)
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

		case conn := <-dispatcher.inboundConns:
			availWorker := <-dispatcher.workerChannel
			availWorker <- conn
			fmt.Println("got a connection from: " + conn.RemoteAddr().String())
		}
	}
}
