package outgoingMsg

// Filter Message Function, Flag to Start Filtering
type OutgoingStopEnforcerMsg struct {
	agent_ip   string
	agent_port string
}

func NewOutgoingStopEnforcerMsg(agent_ip string, agent_port string) OutgoingStopEnforcerMsg {
	return OutgoingStopEnforcerMsg{
		agent_ip:   agent_ip,
		agent_port: agent_port}
}

func (stopEnforcerMsg OutgoingStopEnforcerMsg) String() string {
	return "type: StopEnforcer; payload: Stop Filtering Packets"
}

func (stopEnforcerMsg OutgoingStopEnforcerMsg) GetType() string {
	return "StopEnforcer"
}
