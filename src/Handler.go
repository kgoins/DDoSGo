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

import "visitors"

import "workers"

// import "config"

type Handler struct {
	serverSock        net.Listener
	maxWorkers        int
	maxBuff           int
	dispatcherChannel chan dispatcher.Dispatchable
	dispatcher        *dispatcher.Dispatcher
	alertSystem       *subsystems.AlertSystem
	agentReg          *subsystems.AgentRegistry
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

	listenerSock, sockErr := net.Listen("tcp", port)
	if sockErr != nil {
		fmt.Println("Error binding server socket")
		os.Exit(2)
	}

	agentReg := subsystems.NewAgentRegistry() // Setup agent registry

	monIntval := 12                                                        // TODO: read from conf
	alertSystem := subsystems.NewAlertSystem(agentReg, workers, monIntval) // Setup alert system

	visitors.SetupHandlerVisitors(alertSystem, agentReg)

	return &Handler{
		serverSock:        listenerSock,
		maxWorkers:        workers,
		maxBuff:           buff_size,
		dispatcherChannel: dispatcherChannel,
		dispatcher:        dispatcher,
		alertSystem:       alertSystem,
		agentReg:          agentReg,
	}
}

func (handler *Handler) Run() {
	handler.dispatcher.Run()
	handler.alertSystem.Run()

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

		msgWork := workers.NewMsgDispatchable(conn)
		handler.dispatcherChannel <- msgWork

		chanCap = getMax(chanCap, len(handler.dispatcherChannel))
		fmt.Println("Max Chan cap: ", chanCap)
	}
}

func (handler *Handler) Close() {
	fmt.Println("Closing handler")

	handler.dispatcher.Close()
	handler.alertSystem.Close()

	fmt.Println("closing server sock")
	handler.serverSock.Close()

	close(handler.dispatcherChannel)

	fmt.Println("Handler Closed")
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
