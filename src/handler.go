package main

import "net"
import "fmt"
import "bufio"
import "runtime"
import "strings"

import "config_parsers"

type Handler struct {
	server_sock net.Listener
	max_workers int
}

func NewHandler() Handler {
	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")
	workers := runtime.NumCPU()

	return Handler{server_sock: ln, max_workers: workers}
}

func (handler Handler) Close() {
	handler.server_sock.Close()
}

func (handler Handler) messageDispatcher() {
	// accept connection on port
	conn, _ := handler.server_sock.Accept()

	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		// output message received
		fmt.Print("Message Received:", string(message))

		// sample process for string received
		newmessage := strings.ToUpper(message)

		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Launching server...")

	handler := NewHandler()
	defer handler.Close()
	handler.messageDispatcher()
}
