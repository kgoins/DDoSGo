package main

import "net"
import "fmt"
import "bufio"
import "os"

type Agent struct {
	handlerAddr string
}

func NewAgent() Agent {
	return Agent{handlerAddr: "127.0.0.1:8081"}
}

func (agent Agent) DialHandler() (net.Conn, error) {
	conn, err := net.Dial("tcp", agent.handlerAddr)

	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		return conn, nil
	}
}

func (agent Agent) sendMsg() {
	conn, _ := agent.DialHandler()

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
	agent := NewAgent()
	agent.sendMsg()
}
