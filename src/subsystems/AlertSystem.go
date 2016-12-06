package subsystems

import (
	"dispatcher"
	"fmt"
	"net"
	"outgoingMsg"
	"strings"
	"time"
)

// Alert subsystem structure
type AlertSystem struct {
	agentReg          *AgentRegistry
	dispatcherChannel chan dispatcher.Dispatchable
	dispatcher        *dispatcher.Dispatcher
	monitorIntval     int
	shutdown          chan bool
}

// Init alert subsystem w/ given # of workers, the agent registry to monitor, and the monitor interval
func NewAlertSystem(agentReg *AgentRegistry, workers int, monitorIntval int) *AlertSystem {

	shutdown := make(chan bool)
	dispatcherChannel := make(chan dispatcher.Dispatchable)            // Setup dispatcherChannel non-buffered
	fmt.Println("Setup new alert dispatch with # workers: ", workers)  // Log dispatcher setup
	dispatcher := dispatcher.NewDispatcher(dispatcherChannel, workers) // Create dispatcher

	// Return new AlertSystem ref
	return &AlertSystem{
		agentReg:          agentReg,
		dispatcherChannel: dispatcherChannel,
		dispatcher:        dispatcher,
		monitorIntval:     monitorIntval,
		shutdown:          shutdown,
	}
}

// Close alert subsystem connections
func (alertSystem *AlertSystem) Close() {
	fmt.Println("Closing Alert System")
	alertSystem.shutdown <- true

	alertSystem.dispatcher.Close()
	fmt.Println("Alert System Closed")

	close(alertSystem.dispatcherChannel)
}

// Run alert subsystem
func (alertSystem *AlertSystem) Run() {

	alertSystem.dispatcher.Run()
	alertSystem.MonitorRegistry()
}

// Check agent registry on timed interval (30 seconds) for unflagged agents
func (alertSystem *AlertSystem) MonitorRegistry() {
	go func() {
		for {
			select {
			case <-alertSystem.shutdown:
				fmt.Println("Alert System shutting down")
				return
			default:
				// fmt.Println("Checking Agent Registry For Unresponsive Agents...")
				clean, records := alertSystem.agentReg.CheckRecords() // Check registry for records
				if clean != true {                                    // Records found unresponsive
					fmt.Println("Records Reported Unresponsive")
					for _, record := range records { // Check all unresponsive records
						// fmt.Println(record)
						if record.isFiltering == false { // If agent already filtering, skip msg sending
							filterMsg := outgoingMsg.NewOutgoingFilterMsg(record.agent_ip, record.agent_port) // Send msg to start filtering
							alertSystem.sendFilterMsg(record.agent_ip, record.agent_port, filterMsg)
						} else { // Already filtering, log and ignore
							fmt.Printf("Agent %s%s Already Filtering, Skipping Filter Msg\n", record.agent_ip, record.agent_port)
						}
					}
					// Start alert system
				} else { // No records found unresponsive
					// fmt.Println("No Records Reported Unresponsive")
				}

				time.Sleep(time.Second * time.Duration(alertSystem.monitorIntval)) // Rest for interval length
			}
		}
	}()
}

// Send the filter msg to the agent
func (alertSystem *AlertSystem) sendFilterMsg(agent_ip string, agent_port string, filterMsg outgoingMsg.OutgoingFilterMsg) error {

	fmt.Printf("Sending Alert Msg To Agent\t%s%s\n", agent_ip, agent_port)

	msgBytes, encodeErr := outgoingMsg.EncodeMsg(filterMsg)
	if encodeErr != nil {
		fmt.Println("Error encoding filter fmessage")
		return encodeErr
	}

	conn, err := alertSystem.dialAgentForFiltering(agent_ip, agent_port)

	defer conn.Close()

	if err != nil {
		fmt.Println("err dialing handler", err)
		return err
	}

	fmt.Println("sending message: " + filterMsg.String())
	_, err = conn.Write(msgBytes)
	if err != nil {
		fmt.Println(err)
	}

	return err

}

// Send the stop filter msg to the agent
func (alertSystem *AlertSystem) sendStopEnforcerMsg(agent_ip string, agent_port string, filterMsg outgoingMsg.OutgoingStopEnforcerMsg) error {

	fmt.Printf("Sending Stop Filter Msg To Agent\t%s%s\n", agent_ip, agent_port)

	msgBytes, encodeErr := outgoingMsg.EncodeMsg(filterMsg)
	if encodeErr != nil {
		fmt.Println("Error encoding filter fmessage")
		return encodeErr
	}

	conn, err := alertSystem.dialAgentForFiltering(agent_ip, agent_port)

	defer conn.Close()

	if err != nil {
		fmt.Println("err dialing handler", err)
		return err
	}

	fmt.Println("sending message: " + filterMsg.String())
	_, err = conn.Write(msgBytes)
	if err != nil {
		fmt.Println(err)
	}

	return err

}

// Dial the agent in question for fitlering msgs to begin
func (alertSystem *AlertSystem) dialAgentForFiltering(agent_ip string, agent_port string) (net.Conn, error) {

	conn, err := net.Dial("tcp", agent_ip+agent_port)

	if err != nil {
		fmt.Println(err)
		return conn, err
	} else {
		return conn, nil
	}
}

// Process the data stream data and see if we need to perform an alert
func (alertSystem *AlertSystem) ProcessDataStream(agent_ip string, agent_port string, cpu int, mem int, bytesRecvd int, bytesSent int, isFiltering bool) {

	// fmt.Printf("Processing Data Stream Values\nCPU\t%d\tMEM\t%d\tBytesRecvd\t%d\tBytesSent\t%d\n", cpu, mem, bytesRecvd, bytesSent)
	// If values are strange, alert
	if cpu > 90 {
		fmt.Printf("Cpu Value Anomaly of %d for DataStream of Agent %s%s\n", cpu, agent_ip, agent_port)
		if isFiltering == false { // If not curerntly filtering, send filter msg
			filterMsg := outgoingMsg.NewOutgoingFilterMsg(agent_ip, agent_port)
			alertSystem.sendFilterMsg(agent_ip, agent_port, filterMsg)
			alertSystem.agentReg.SetAgentAsFiltering(agent_ip, agent_port)

			//Send filtering msg to Agent's trace
			trace := alertSystem.agentReg.ReturnTrace(agent_ip, agent_port)

			for _, ip := range trace {
				ipPort := strings.Split(ip, ":")
				ipPort[1] = ":" + ipPort[1]

				//Check if trace is in the registry
				reg := *alertSystem.agentReg
				_, exists := reg.registry[ip]
				if exists {
					traceFilterMsg := outgoingMsg.NewOutgoingFilterMsg(ipPort[0], ipPort[1], true)
					alertSystem.sendFilterMsg(ipPort[0], ipPort[1], traceFilterMsg)
					alertSystem.agentReg.SetAgentAsFiltering(ipPort[0], ipPort[1])
				}
			}

		} else { // Already filtering, log and ignore
			fmt.Println("Ignoring Anomaly As Agent Already Filtering")
		}

	} else if (cpu <= 90) && (isFiltering == true) {

		fmt.Printf("Stopping Filtering For Agent %s%s\n", agent_ip, agent_port)

		stopEnforcerMsg := outgoingMsg.NewOutgoingStopEnforcerMsg(agent_ip, agent_port) // Create msg to stop filter
		alertSystem.sendStopEnforcerMsg(agent_ip, agent_port, stopEnforcerMsg)          // Send new filter msg
		alertSystem.agentReg.ClearAgentFilteringStatus(agent_ip, agent_port)            // Clear agent record filtering status

	} else {
		// fmt.Println("No DataStream Anomalies Detected")
	}
	// Else all clear
}
