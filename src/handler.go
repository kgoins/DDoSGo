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
	workers := runtime.NumCPU() * 20
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

	var chanCap int = 0

	for {
		conn, _ := handler.server_sock.Accept()
		handler.msgChannel <- conn

		chanCap = getMax(chanCap, len(handler.msgChannel))
		fmt.Println("Max Chan cap: ", chanCap)
	}
}

func getMax(num1, num2 int) int {
	if num1 > num2 {
		return num1
	} else {
		return num2
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	handler := NewHandler()
	handler.Run()
}
