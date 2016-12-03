package subsystems

import (
	"dispatcher"
	"fmt"
)

// Alert subsystem structure
type AlertSystem struct {
	dispatcherChannel chan dispatcher.Dispatchable
	dispatcher        *dispatcher.Dispatcher
}

// Init alert subsystem w/ given # of workers
func NewAlertSystem(workers int) *AlertSystem {

	dispatcherChannel := make(chan dispatcher.Dispatchable)            // Setup dispatcherChannel non-buffered
	fmt.Println("Setup new alert dispatch with # workers: ", workers)  // Log dispatcher setup
	dispatcher := dispatcher.NewDispatcher(dispatcherChannel, workers) // Create dispatcher

	// Return new AlertSystem ref
	return &AlertSystem{
		dispatcherChannel: dispatcherChannel,
		dispatcher:        dispatcher,
	}
}

// Close alert subsystem connections
func (alertSystem *AlertSystem) Close() {
	fmt.Println("Closing Alert System Connections...")

	// Close dispatcher connections & channel
	alertSystem.dispatcher.Close()
	close(alertSystem.dispatcherChannel)

}

// Run alert subsystem
func (alertSystem *AlertSystem) Run() {
	fmt.Println("Alert System Starting...")
	alertSystem.dispatcher.Run()
}
