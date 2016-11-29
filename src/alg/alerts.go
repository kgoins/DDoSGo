package alg

// Alert subsystem structure
type AlertSystem struct {
	dispatcher *AlertDispatcher
}

// Init alert subsystem
func NewAlertSystem() *AlertSystem {
	return &AlertSystem{}
}
