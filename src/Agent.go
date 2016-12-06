package main

import "net"
import "fmt"
import "errors"

import "os"
import "syscall"
import "os/signal"

import "outgoingMsg"
import "config"

import "dispatcher"
import "subsystems"

import "workers"
import "visitors"

type Agent struct {
	handlerAddr        string
	serverSock         net.Listener
	collectionInterval int
	msgChannel         chan outgoingMsg.OutgoingMsg
	shutdown           chan bool

	agent_port   string
	agent_ip     string
	handler_ip   string
	handler_port string
	trace        []string

	collector *subsystems.DataCollector
	enforcer  *subsystems.Enforcer

	dispatcherChannel chan dispatcher.Dispatchable
	dispatcher        *dispatcher.Dispatcher
}

func NewAgent(portNum string, nfQueue_num string) (Agent, error) {
	agentConf, err := config.ReadAgentConf()

	var handlerAddr string = agentConf.HandlerAddr + ":" + agentConf.HandlerPort
	fmt.Println("Connecting to handler: " + handlerAddr)

	aIp, _ := getIP()
	hIp := agentConf.HandlerAddr
	hPort := agentConf.HandlerPort
	tracert := agentConf.Traceroute

	dispatcherChannel := make(chan dispatcher.Dispatchable)
	msgChannel := make(chan outgoingMsg.OutgoingMsg)
	shutdown := make(chan bool)

	port := portNum
	port = ":" + port

	collectionInterval := 2
	sendInterval := 5
	collector := subsystems.NewDataCollector(aIp, port, hIp, hPort, msgChannel, collectionInterval, sendInterval)

	numWorkers := 2 // TODO: read from conf
	dispatcher := dispatcher.NewDispatcher(dispatcherChannel, numWorkers)

	enforcer := subsystems.NewEnforcer(nfQueue_num)

	visitors.SetupAgentVisitors(enforcer) // Setup enforcement visitor
	serverSock, _ := net.Listen("tcp", port)

	return Agent{handlerAddr: handlerAddr,
		serverSock: serverSock,

		collectionInterval: collectionInterval,
		collector:          collector,
		dispatcher:         dispatcher,
		enforcer:           enforcer,

		dispatcherChannel: dispatcherChannel,
		shutdown:          shutdown,
		msgChannel:        msgChannel,

		agent_port:   port,
		agent_ip:     aIp,
		handler_ip:   hIp,
		handler_port: hPort,
		trace:        tracert}, err
}

func (agent Agent) Start() {
	agent.signalHandler()

	//Build and send register msg
	regMsg := outgoingMsg.NewOutgoingRegisterMsg(agent.agent_ip, agent.handler_ip, agent.handler_port, agent.trace, agent.agent_port, false)
	agent.sendMsg(regMsg)

	go agent.msgSender()

	// Start subsystems
	agent.collector.Start()
	agent.dispatcher.Run()

	agent.msgReceiver()

	fmt.Println("Msg Receiver died, waiting on Close")
	<-agent.shutdown
}

func (agent Agent) Close() {
	fmt.Println("Closing subsystems")
	agent.collector.Close()
	agent.dispatcher.Close()
	agent.enforcer.Close()
	fmt.Println("All Subsystems Closed")

	closeRegMsg := outgoingMsg.NewOutgoingRegisterMsg(agent.agent_ip, agent.handler_ip, agent.handler_port, agent.trace, agent.agent_port, true)
	agent.sendMsg(closeRegMsg)
	fmt.Println("Sent shutdown Msg to Handler")

	agent.shutdown <- true

	agent.serverSock.Close()
	fmt.Println("Server Sock Closed")

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
			fmt.Println("Msg Sender Closing")
			return

		case msg := <-agent.msgChannel:
			sendErr := agent.sendMsg(msg)
			if sendErr != nil {
				fmt.Println(sendErr)
			}
		}
	}
}

func (agent Agent) sendMsg(msg outgoingMsg.OutgoingMsg) error {
	msgBytes, encodeErr := outgoingMsg.EncodeMsg(msg)
	if encodeErr != nil {
		fmt.Println("Error encoding message")
		return encodeErr
	}

	conn, err := agent.dialHandler()
	defer conn.Close()
	if err != nil {
		fmt.Println("err dialing handler", err)
		return err
	}

	fmt.Println("sending message: " + msg.String())
	_, err = conn.Write(msgBytes)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (agent Agent) msgReceiver() {
	for {
		conn, err := agent.serverSock.Accept()
		if err != nil {
			agent.ntwkErrHandler(err)
			return
		}

		msgWork := workers.NewMsgDispatchable(conn)
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
		fmt.Println("Unknown network error:", err)
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

func getIP() (string, error) {
	ipStr := ""

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		return ipStr, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return ipStr, errors.New("Unable to obtain agent IP")
}

// *** AGENT'S MAIN *** //
func main() {

	//Get command line arguments
	args := os.Args
	if len(args) != 3 {
		fmt.Printf("USAGE: ./agent.go port_num nfQueue_num\nThe Agent takes a port number and Iptables netfilter Queue number as arguments")
		return
	}

	agent, newAgentErr := NewAgent(args[1], args[2])
	if newAgentErr != nil {
		fmt.Println("error creating agent")
		fmt.Println(newAgentErr)
		return
	}
	agent.Start()
}
