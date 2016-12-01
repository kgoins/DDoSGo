package main

import "net"
import "fmt"
import "runtime"

import "os"
import "syscall"
import "os/signal"

import "dispatcher"

// import "config"

type Handler struct {
	serverSock        net.Listener
	maxWorkers        int
	maxBuff           int
	dispatcherChannel chan dispatcher.Dispatchable
	killsig           chan bool
	dispatcher        *dispatcher.Dispatcher
}

func NewHandler() *Handler {
	port := ":1337" // TODO: read from conf

	buff_size := 1000 // TODO: read from conf
	dispatcherChannel := make(chan dispatcher.Dispatchable, buff_size)

	workers := runtime.NumCPU() * 10 // TODO: read from conf
	dispatcher := dispatcher.NewDispatcher(dispatcherChannel, workers)

	fmt.Println("Num workers: ", workers)
	listenerSock, _ := net.Listen("tcp", port)

	return &Handler{
		serverSock:        listenerSock,
		maxWorkers:        workers,
		maxBuff:           buff_size,
		dispatcherChannel: dispatcherChannel,
		killsig:           make(chan bool),
		dispatcher:        dispatcher}
}

func (handler *Handler) Run() {
	fmt.Println("Starting Handler...")

	handler.dispatcher.Run()

	handler.signalHandler()
	handler.serve()
}

func (handler *Handler) serve() {
	chanCap := 0
	for {
		conn, err := handler.serverSock.Accept()
		if err != nil {
			handler.serverErrHandler(err)
		}

		msgWork := dispatcher.NewMsgDispatchable(conn)
		handler.dispatcherChannel <- msgWork

		chanCap = getMax(chanCap, len(handler.dispatcherChannel))
		fmt.Println("Max Chan cap: ", chanCap)
	}
}

func (handler *Handler) Close() {
	handler.dispatcher.Close()
	handler.serverSock.Close()

	close(handler.dispatcherChannel)
	close(handler.killsig)

	fmt.Println("Closing handler")
	os.Exit(1)
}

func (handler *Handler) serverErrHandler(err error) {
	switch errType := err.(type) {
	case *net.OpError:
		if errType.Op == "accept" {
			println("Server socket closed, shutting down")
			handler.Close()
		}

	default:
		fmt.Println(err)
	}
}

func (handler *Handler) signalHandler() {
	osKillsig := make(chan os.Signal, 1)
	signal.Notify(osKillsig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-osKillsig
		handler.Close()
	}()
}

// Util functions
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
