package main

import "net"
import "fmt"
import "runtime"
import "log"

import "os"
import "syscall"
import "os/signal"

import "dispatcher"
import "subsystems"

// import "config"

type Handler struct {
	serverSock        net.Listener
	maxWorkers        int
	maxBuff           int
	dispatcherChannel chan dispatcher.Dispatchable
	dispatcher        *dispatcher.Dispatcher
}

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func NewHandler() *Handler {
	port := ":1337" // TODO: read from conf

	buff_size := 1000 // TODO: read from conf
	dispatcherChannel := make(chan dispatcher.Dispatchable, buff_size)

	workers := runtime.NumCPU() * 10 // TODO: read from conf
	fmt.Println("Num workers: ", workers)
	dispatcher := dispatcher.NewDispatcher(dispatcherChannel, workers)

	listenerSock, _ := net.Listen("tcp", port)

	return &Handler{
		serverSock:        listenerSock,
		maxWorkers:        workers,
		maxBuff:           buff_size,
		dispatcherChannel: dispatcherChannel,
		dispatcher:        dispatcher}
}

func (handler *Handler) Run() {
	fmt.Println("Starting Handler...")

	handler.dispatcher.Run()

	handler.signalHandler()

	//Intitalize the Agent Registry
	fmt.Println("Starting Agent Registry...")
	network.Start()

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
		handler.Close()
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

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func setupLogger() {

	direxists, err := exists("./logs")
	if !direxists {
		direrr := os.Mkdir("./logs", 0766)
		if direrr != nil {
			fmt.Printf("Error init logging dir: %v", direrr)
			os.Exit(1)
		}
	}

	fexists, err := exists("./logs/handler_log.txt")
	if !fexists {
		os.Create("./logs/handler_log.txt")
	}

	errlog, err := os.Open("./logs/handler_log.txt")
	if err != nil {
		fmt.Printf("Error initalizing error logging: %v", err)
		os.Exit(1)
	}

	Trace = log.New(errlog, "Application Log: ", log.Lshortfile|log.LstdFlags)

}

// *** MAIN *** //
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setupLogger()
	handler := NewHandler()
	handler.Run()
}
