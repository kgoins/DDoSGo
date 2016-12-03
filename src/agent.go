package main

import "net"
import "fmt"

import "os"
import "syscall"
import "os/signal"

import "msgs"
import "config"

import "dispatcher"
import "subsystems"

type Agent struct {
	handlerAddr        string
	serverSock         net.Listener
	collectionInterval int
	msgChannel         chan msgs.Msg
	shutdown           chan bool

	collector subsystems.DataCollector

	dispatcherChannel chan dispatcher.Dispatchable
	dispatcher        *dispatcher.Dispatcher
}

func NewAgent() (Agent, error) {
	agentConf, err := config.ReadAgentConf()

	var handlerAddr string = agentConf.HandlerAddr + ":" + agentConf.HandlerPort
	fmt.Println("Connecting to handler: " + handlerAddr)

	dispatcherChannel := make(chan dispatcher.Dispatchable)
	msgChannel := make(chan msgs.Msg)
	shutdown := make(chan bool)

	collectionInterval := 2
	sendInterval := 5
	collector := subsystems.NewDataCollector(msgChannel, collectionInterval, sendInterval)

	numWorkers := 2 // TODO: read from conf
	dispatcher := dispatcher.NewDispatcher(dispatcherChannel, numWorkers)

	port := ":1338" // TODO: read from conf
	serverSock, _ := net.Listen("tcp", port)

	return Agent{handlerAddr: handlerAddr,
		serverSock:         serverSock,
		collectionInterval: collectionInterval,
		collector:          collector,
		dispatcher:         dispatcher,
		dispatcherChannel:  dispatcherChannel,
		shutdown:           shutdown,
		msgChannel:         msgChannel}, err
}

func (agent Agent) Start() {
	agent.signalHandler()

	agent.collector.Start()
	agent.dispatcher.Run()

	go agent.msgSender()
	agent.msgReceiver()
}

func (agent Agent) Close() {
	agent.shutdown <- true

	agent.serverSock.Close()
	agent.collector.Close()
	agent.dispatcher.Close()

	close(agent.msgChannel)
	close(agent.shutdown)

	fmt.Println("agent closed")
	os.Exit(1)
}

func (agent Agent) dialHandler() (net.Conn, error) {
	conn, err := net.Dial("tcp", agent.handlerAddr)

	if err != nil {
		fmt.Println(err)
		return conn, err
	} else {
		return conn, nil
	}
}

func (agent Agent) msgSender() {
	for {
		select {
		case <-agent.shutdown:
			return

		case msg := <-agent.msgChannel:
			conn, dialErr := agent.dialHandler()
			if dialErr != nil {
				agent.ntwkErrHandler(dialErr)
			}

			msgData := msgs.EncodeMsg(msg)

			fmt.Println("sending message: " + msg.String())
			_, writeErr := conn.Write(msgData)
			if writeErr != nil {
				agent.ntwkErrHandler(writeErr)
			}

			conn.Close()
		}
	}
}

func (agent Agent) msgReceiver() {
	for {
		conn, err := agent.serverSock.Accept()
		if err != nil {
			agent.ntwkErrHandler(err)
		}

		msgWork := dispatcher.NewMsgDispatchable(conn)
		agent.dispatcherChannel <- msgWork
	}
}

func (agent Agent) ntwkErrHandler(err error) {
	switch errType := err.(type) {
	case *net.OpError:
		if errType.Op == "accept" {
			println("Server socket closed, shutting down")
			// agent.Close()
		}

	default:
		fmt.Println(err)
		agent.Close()
	}
}

func (agent Agent) signalHandler() {
	osKillsig := make(chan os.Signal, 1)
	signal.Notify(osKillsig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-osKillsig
		fmt.Println("OS SIGTERM received, agent shutting down")
		agent.Close()
	}()
}

// *** AGENT'S MAIN *** //
func main() {
	agent, newAgentErr := NewAgent()
	if newAgentErr != nil {
		fmt.Println("error creating agent")
		fmt.Println(newAgentErr)
		return
	}
	agent.Start()
}
