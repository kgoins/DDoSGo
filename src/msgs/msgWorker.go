package msgs

import "net"
import "fmt"
import "bufio"

import "time"

type MsgWorker struct {
	workerPool  chan chan net.Conn
	connChannel chan net.Conn
	quit        chan bool
}

func NewMsgWorker(workerPool chan chan net.Conn) MsgWorker {
	connChannel := make(chan net.Conn)
	quit := make(chan bool)

	return MsgWorker{
		workerPool:  workerPool,
		connChannel: connChannel,
		quit:        quit}
}

func (worker *MsgWorker) Start() {
	go func() {
		for {
			worker.workerPool <- worker.connChannel

			select {
			case <-worker.quit:
				return

			case conn := <-worker.connChannel:
				handleConn(conn)
			}
		}
	}()
}

func handleConn(conn net.Conn) {
	fmt.Println("We have work!")
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(message)
	conn.Close()
	time.Sleep(100 * time.Millisecond)
}
