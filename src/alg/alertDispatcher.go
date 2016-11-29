package alg

// Alert dispatcher struct
type AlertDispatcher struct {
	maxWorkers int
}

// Create new alert dispatcher
func NewAlertDispatcher(workers int) *AlertDispatcher {
	return &AlertDispatcher{
		maxWorkers: workers
	}
}

// For given stream decide if respective agent under attack
func (alertDispatcher *AlertDispatcher) EvaluateAgentThreat() bool {

}

// Generate alert chain
func GenerateAlertChain() {

}

// Send filtering msg
func SendFilteringMsg() {

}
