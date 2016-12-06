package outgoingMsg

// Filter Message Function, Flag to Start Filtering
type OutgoingFilterMsg struct {
	agent_ip           string
	agent_port         string
	startFilterPackets bool
}

func NewOutgoingFilterMsg(agent_ip string, agent_port string, startFilterPackets bool) OutgoingFilterMsg {
	return OutgoingFilterMsg{
		agent_ip:           agent_ip,
		agent_port:         agent_port,
		startFilterPackets: startFilterPackets}
}

func (filterMsg OutgoingFilterMsg) String() string {
	return "type: Filter; payload: Start / stop filtering packets"
}

func (filterMsg OutgoingFilterMsg) GetType() string {
	return "Filter"
}
