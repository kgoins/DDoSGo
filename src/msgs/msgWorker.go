package msgs

import "net"
import "fmt"
import "sync"
import "encoding/json"

import "time"

type MsgWorker struct {
	workerPool  chan chan net.Conn
	connChannel chan net.Conn
	quit        chan bool
	waitgroup   *sync.WaitGroup
}

func NewMsgWorker(workerPool chan chan net.Conn) *MsgWorker {
	connChannel := make(chan net.Conn)
	quit := make(chan bool)

	return &MsgWorker{
		workerPool:  workerPool,
		connChannel: connChannel,
		waitgroup:   &sync.WaitGroup{},
		quit:        quit}
}

func (worker *MsgWorker) Close() {
	worker.quit <- true
	worker.waitgroup.Wait()
}

func (worker *MsgWorker) Start() {
	worker.waitgroup.Add(1)
	go func() {
		defer worker.waitgroup.Done()
		for {
			select {
			case <-worker.quit:
				return

			case worker.workerPool <- worker.connChannel:
				conn, closed := <-worker.connChannel
				if !closed {
					fmt.Println("worker pool closed")
					return
				}

				handleConn(conn)
			}
		}
	}()
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)

	var msg Msg
	err := decoder.Decode(&msg)

	fmt.Println(msg, err)

	time.Sleep(100 * time.Millisecond)
}
