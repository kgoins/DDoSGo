package msgs

import "net"
import "fmt"

type MsgDispatcher struct {
	maxWorkers    int
	inboundConns  chan net.Conn
	workerChannel chan chan net.Conn
	workerPool    []*MsgWorker
	shutdown      chan bool
}

func NewMsgDispatcher(inboundConns chan net.Conn, workers int) *MsgDispatcher {
	workerChannel := make(chan chan net.Conn)
	shutdown := make(chan bool)

	return &MsgDispatcher{
		maxWorkers:    workers,
		inboundConns:  inboundConns,
		workerChannel: workerChannel,
		workerPool:    make([]*MsgWorker, workers),
		shutdown:      shutdown}
}

func (dispatcher *MsgDispatcher) Close() {
	// kill dispatcher goroutine
	dispatcher.shutdown <- true

	for _, worker := range dispatcher.workerPool {
		worker.Close()
	}

	close(dispatcher.workerChannel)
	fmt.Println("all workers shutdown")
}

func (dispatcher *MsgDispatcher) Run() {

	for i := 0; i < dispatcher.maxWorkers; i++ {
		worker := NewMsgWorker(dispatcher.workerChannel)
		dispatcher.workerPool[i] = worker
		worker.Start()
	}

	go dispatcher.dispatch()
}

func (dispatcher *MsgDispatcher) dispatch() {
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
