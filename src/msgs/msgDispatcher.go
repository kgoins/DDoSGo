package msgs

import "net"

type MsgDispatcher struct {
	maxWorkers   int
	inboundConns chan net.Conn
	workerPool   chan chan net.Conn
}

func NewMsgDispatcher(inboundConns chan net.Conn, workers int) *MsgDispatcher {
	workerPool := make(chan chan net.Conn)

	return &MsgDispatcher{
		maxWorkers:   workers,
		inboundConns: inboundConns,
		workerPool:   workerPool}
}

func (dispatcher *MsgDispatcher) Run() {
	for i := 0; i < dispatcher.maxWorkers; i++ {
		worker := NewMsgWorker(dispatcher.workerPool)
		worker.Start()
	}

	go dispatcher.dispatch()
}

func (dispatcher *MsgDispatcher) dispatch() {
	for {
		select {
		case conn := <-dispatcher.inboundConns:
			connChannel := <-dispatcher.workerPool
			connChannel <- conn
		}
	}
}
