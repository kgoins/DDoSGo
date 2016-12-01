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
	max_workers       int
	max_buff          int
	dispatcherChannel chan dispatcher.Dispatchable
	killsig           chan bool
	dispatcher        *dispatcher.Dispatcher
}

func NewHandler() *Handler {
	port := ":1337"

	buff_size := 1000 // TODO: read from conf
	dispatcherChannel := make(chan dispatcher.Dispatchable, buff_size)

	workers := runtime.NumCPU() * 10 // TODO: read from conf
	dispatcher := dispatcher.NewDispatcher(dispatcherChannel, workers)

	fmt.Println("Num workers: ", workers)
	listenerSock, _ := net.Listen("tcp", port)

	return &Handler{
		serverSock:        listenerSock,
		max_workers:       workers,
		max_buff:          buff_size,
		dispatcherChannel: dispatcherChannel,
		killsig:           make(chan bool),
		dispatcher:        dispatcher}
}

func (handler *Handler) Run() {
	fmt.Println("Starting Handler...")

	handler.dispatcher.Run()

	handler.signalHandler()
	handler.serve()

	<-handler.killsig
	fmt.Println("Handler closed")
}

func (handler *Handler) serve() {
	chanCap := 0
	for {
		conn, err := handler.serverSock.Accept()
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

		msgWork := dispatcher.NewMsgDispatchable(conn)
		handler.dispatcherChannel <- msgWork

		chanCap = getMax(chanCap, len(handler.dispatcherChannel))
		fmt.Println("Max Chan cap: ", chanCap)
	}
}

func (handler *Handler) Close() {
	fmt.Println("Closing handler")

	handler.dispatcher.Close()
	handler.serverSock.Close()

	close(handler.dispatcherChannel)
	close(handler.killsig)

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
