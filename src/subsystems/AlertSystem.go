package subsystems

import (
	"dispatcher"
	"fmt"
	"outgoingMsg"
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
				fmt.Println("Checking Agent Registry For Unresponsive Agents...")
				clean, records := alertSystem.agentReg.CheckRecords() // Check registry for records
				if clean != true {                                    // Records found unresponsive
					fmt.Println("Records Reported Unresponsive")
					for _, record := range records {
						fmt.Println(record)
						filterMsg := outgoingMsg.NewOutgoingFilterMsg(record.agent_ip, true) // Send msg to start filtering
						alertSystem.sendAlertMsg(record.agent_ip, filterMsg)
					}
					// Start alert system
				} else { // No records found unresponsive
					fmt.Println("No Records Reported Unresponsive")
				}

				time.Sleep(time.Second * time.Duration(alertSystem.monitorIntval)) // Rest for interval length
			}
		}
	}()
}

// Send the alert msg
func (alertSystem *AlertSystem) sendAlertMsg(agent_ip string, filterMsg outgoingMsg.OutgoingFilterMsg) {
	fmt.Printf("Sending Alert Msg\nIP\t%s\n", agent_ip)
}

// Process the data stream data and see if we need to perform an alert
func (alertSystem *AlertSystem) ProcessDataStream(agent_ip string, cpu int, mem int, bytesRecvd int, bytesSent int) {

	// fmt.Printf("Processing Data Stream Values\nCPU\t%d\tMEM\t%d\tBytesRecvd\t%d\tBytesSent\t%d\n", cpu, mem, bytesRecvd, bytesSent)
	// If values are strange, alert
	if cpu > 5 {
		fmt.Printf("Cpu Value Anomaly of %d for DataStream from Agent %s\n", cpu, agent_ip)
		filterMsg := outgoingMsg.NewOutgoingFilterMsg(agent_ip, true)
		alertSystem.sendAlertMsg(agent_ip, filterMsg)
	} else {
		fmt.Println("No DataStream Anomalies Detected")
	}
	// Else all clear
}
