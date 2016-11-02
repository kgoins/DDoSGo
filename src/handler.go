package main

import "net"
import "fmt"
import "runtime"
import "msgs"

// import "config"

type Handler struct {
	server_sock net.Listener
	max_workers int
	max_buff    int
	msgChannel  chan net.Conn
}

func NewHandler() *Handler {
	port := ":1337"
	workers := runtime.NumCPU() * 1000
	fmt.Println("Num workers: ", workers)
	buff_size := 1000

	// listen on all interfaces
	listenerSock, _ := net.Listen("tcp", port)

	return &Handler{
		server_sock: listenerSock,
		max_workers: workers,
		max_buff:    buff_size,
		msgChannel:  make(chan net.Conn, buff_size)}
}

func (handler *Handler) Close() {
	handler.server_sock.Close()
}

func (handler *Handler) Run() {
	defer handler.Close()
	fmt.Println("Starting Handler...")

	dispatcher := msgs.NewMsgDispatcher(handler.msgChannel, handler.max_workers)
	dispatcher.Run()

	for {
		conn, _ := handler.server_sock.Accept()
		handler.msgChannel <- conn
		fmt.Println("Chan cap: ", len(handler.msgChannel))
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	handler := NewHandler()
	handler.Run()
}
