package main

import "net"
import "fmt"
import "runtime"

import "os"
import "syscall"
import "os/signal"

import "msgs"
import "network"

// import "config"

type Handler struct {
	server_sock net.Listener
	max_workers int
	max_buff    int
	msgChannel  chan net.Conn
	killsig     chan bool
	dispatcher  *msgs.MsgDispatcher
}

func NewHandler() *Handler {
	port := ":1337"
	workers := runtime.NumCPU() * 10
	fmt.Println("Num workers: ", workers)
	buff_size := 1000

	// listen on all interfaces
	msgChannel := make(chan net.Conn, buff_size)
	listenerSock, _ := net.Listen("tcp", port)
	dispatcher := msgs.NewMsgDispatcher(msgChannel, workers)

	return &Handler{
		server_sock: listenerSock,
		max_workers: workers,
		max_buff:    buff_size,
		msgChannel:  msgChannel,
		killsig:     make(chan bool),
		dispatcher:  dispatcher}
}

func (handler *Handler) Run() {
	fmt.Println("Starting Handler...")
	handler.dispatcher.Run()

	handler.signalHandler()

	//Intitalize the Agent Registry
	fmt.Println("Starting Agent Registry...")
	network.Start()

	handler.serve()
	<-handler.killsig
	fmt.Println("Handler closed")
}

func (handler *Handler) serve() {
	chanCap := 0
	for {
		conn, err := handler.server_sock.Accept()
		if err != nil {
			switch errType := err.(type) {
			case *net.OpError:
				if errType.Op == "accept" {
					println("Server socket closed")
					return
				}

			default:
				fmt.Println(err)
			}
		}

		handler.msgChannel <- conn

		chanCap = getMax(chanCap, len(handler.msgChannel))
		fmt.Println("Max Chan cap: ", chanCap)
	}
}

func (handler *Handler) Close() {
	fmt.Println("Closing handler")
	handler.dispatcher.Close()

	handler.server_sock.Close()
	handler.killsig <- true
}

func (handler *Handler) signalHandler() {
	killsig := make(chan os.Signal, 1)
	signal.Notify(killsig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-killsig
		handler.Close()
		os.Exit(1)
	}()
}

func getMax(num1, num2 int) int {
	if num1 > num2 {
		return num1
	} else {
		return num2
	}
}

// *** MAIN *** //
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	handler := NewHandler()
	handler.Run()
}
