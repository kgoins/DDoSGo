package main

import "net"
import "fmt"

// import "bufio"
// import "os"

import "msgs"
import "config"

type Agent struct {
	handlerAddr string
}

func NewAgent() (Agent, error) {
	agentConf, err := config.ReadAgentConf()

	var handlerAddr string = agentConf.HandlerAddr + ":" + agentConf.HandlerPort
	fmt.Println(handlerAddr)

	return Agent{handlerAddr: handlerAddr}, err
}

func (agent Agent) DialHandler() (net.Conn, error) {
	conn, err := net.Dial("tcp", agent.handlerAddr)

	if err != nil {
		fmt.Println(err)
		return conn, err
	} else {
		return conn, nil
	}
}

func (agent Agent) sendMsg() error {
	conn, err := agent.DialHandler()
	if err != nil {
		return err
	}
	defer conn.Close()

	msg := msgs.NewDebugMsg("Hello World!")
	msgData := msgs.EncodeMsg(msg)

	fmt.Println(string(msgData))

	// send to socket
	conn.Write(msgData)
	fmt.Println("sending message: " + msg.String())

	return nil
}

func main() {
	agent, newAgentErr := NewAgent()
	if newAgentErr != nil {
		fmt.Println("error creating agent")
		fmt.Println(newAgentErr)
		return
	}
	agent.sendMsg()
}
