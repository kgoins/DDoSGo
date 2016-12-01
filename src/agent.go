package main

import "net"
import "fmt"

import "msgs"
import "config"

import "data"

type Agent struct {
	handlerAddr        string
	msgChannel         chan Msg
	serverSock         net.Listener
	collectionInterval int
}

func NewAgent() (Agent, error) {
	agentConf, err := config.ReadAgentConf()

	var handlerAddr string = agentConf.HandlerAddr + ":" + agentConf.HandlerPort
	fmt.Println("Connecting to handler: " + handlerAddr)

	collectionInterval := 15

	msgChannel := make(chan Msg)

	port := ":1338" // TODO: read from conf
	serverSock, _ := net.Listen("tcp", port)

	return Agent{handlerAddr: handlerAddr,
		serverSock:         listenerSock,
		collectionInterval: collectionInterval,
		msgChannel:         msgChannel}, err
}

func (agent Agent) Start() {
	go msgSender()

	collector := data.NewDataCollector(agent.msgChannel, agent.collectionInterval)
	collector.Start()

	// msgReceiver()
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
		case <-shutdown:
			return

		case msg := <-agent.msgChannel:
			conn, err := agent.dialHandler()
			if err != nil {
				return err
			}

			msgData := msgs.EncodeMsg(msg)

			conn.Write(msgData)
			fmt.Println("sending message: " + msg.String())

			conn.Close()
		}
	}
}

func (agent Agent) msgReceiver() {
}

func main() {
	agent, newAgentErr := NewAgent()
	if newAgentErr != nil {
		fmt.Println("error creating agent")
		fmt.Println(newAgentErr)
		return
	}
	agent.Start()
}
