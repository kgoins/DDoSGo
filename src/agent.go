package main

import "net"
import "fmt"
import "bufio"
import "os"

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
		os.Exit(1)
		return nil, err
	} else {
		return conn, nil
	}
}

func (agent Agent) sendMsg() {
	conn, _ := agent.DialHandler()
	defer conn.Close()

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')

		// send to socket
		fmt.Fprintf(conn, text+"\n")

		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
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
