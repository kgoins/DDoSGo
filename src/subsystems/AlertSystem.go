package subsystems

import (
	"dispatcher"
	"fmt"
	"outgoingMsg"
	"time"
)

// Alert subsystem structure
type AlertSystem struct {
	dispatcherChannel chan dispatcher.Dispatchable
	dispatcher        *dispatcher.Dispatcher

	monitorIntval int
	shutdown      chan bool
}

// Init alert subsystem w/ given # of workers
func NewAlertSystem(workers int, monitorIntval int) *AlertSystem {

	shutdown := make(chan bool)
	dispatcherChannel := make(chan dispatcher.Dispatchable)            // Setup dispatcherChannel non-buffered
	fmt.Println("Setup new alert dispatch with # workers: ", workers)  // Log dispatcher setup
	dispatcher := dispatcher.NewDispatcher(dispatcherChannel, workers) // Create dispatcher

	// Return new AlertSystem ref
	return &AlertSystem{
		dispatcherChannel: dispatcherChannel,
		dispatcher:        dispatcher,
		monitorIntval:     monitorIntval,
		shutdown:          shutdown,
	}
}

// Close alert subsystem connections
func (alertSystem *AlertSystem) Close() {
	fmt.Println("Closing Alert System Connections...")

	// Close dispatcher connections & channel
	alertSystem.shutdown <- true
	alertSystem.dispatcher.Close()
	close(alertSystem.dispatcherChannel)

}

// Run alert subsystem
func (alertSystem *AlertSystem) Run() {
	fmt.Println("Alert System Starting...")
	alertSystem.MonitorRegistry()
	alertSystem.dispatcher.Run()
}

// Check agent registry on timed interval (30 seconds) for unflagged agents
func (alertSystem *AlertSystem) MonitorRegistry() {
	go func() {
		for {
			select {
			case <-alertSystem.shutdown:
				return
			default:
				fmt.Println("Checking Agent Registry For Unresponsive Agents...")

				time.Sleep(time.Second * time.Duration(alertSystem.monitorIntval))
			}
		}
	}()
}

// For received DataStream evaluate if agent is under attack
func (alertSystem *AlertSystem) ProcessDataStream(stream outgoingMsg.OutgoingDataStream) {

}
