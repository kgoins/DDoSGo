package subsystems

import "dispatcher"

// Alert subsystem structure
type AlertSystem struct {
	dispatcherChannel chan dispatcher.Dispatchable
	dispatcher        *dispatcher.Dispatcher
	monitorIntval     int
	shutdown          chan bool
}
