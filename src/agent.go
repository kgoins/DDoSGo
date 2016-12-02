package main

import "net"
import "fmt"

import "os"
import "syscall"
import "os/signal"

import "msgs"
import "config"

import "data"

type Agent struct {
	handlerAddr        string
	serverSock         net.Listener
	collectionInterval int
	msgChannel         chan msgs.Msg
	shutdown           chan bool

	collector data.DataCollector
}

func NewAgent() (Agent, error) {
	agentConf, err := config.ReadAgentConf()

	var handlerAddr string = agentConf.HandlerAddr + ":" + agentConf.HandlerPort
	fmt.Println("Connecting to handler: " + handlerAddr)

	msgChannel := make(chan msgs.Msg)
	shutdown := make(chan bool)

	collectionInterval := 2
	sendInterval := 5
	collector := data.NewDataCollector(msgChannel, collectionInterval, sendInterval)

	port := ":1338" // TODO: read from conf
	serverSock, _ := net.Listen("tcp", port)

	return Agent{handlerAddr: handlerAddr,
		serverSock:         serverSock,
		collectionInterval: collectionInterval,
		collector:          collector,
		shutdown:           shutdown,
		msgChannel:         msgChannel}, err
}

func (agent Agent) Start() {
	agent.signalHandler()

	agent.collector.Start()

	agent.msgSender()
	// msgReceiver()
}

func (agent Agent) Close() {
	agent.shutdown <- true

	agent.collector.Close()
	agent.serverSock.Close()

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
}

func (agent Agent) ntwkErrHandler(err error) {
	switch errType := err.(type) {
	case *net.OpError:
		if errType.Op == "accept" {
			println("Server socket closed, shutting down")
			agent.Close()
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

// *** MAIN *** //
func main() {
	agent, newAgentErr := NewAgent()
	if newAgentErr != nil {
		fmt.Println("error creating agent")
		fmt.Println(newAgentErr)
		return
	}
	agent.Start()
}
