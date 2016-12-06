package outgoingMsg

// Filter Message Function, Flag to Start Filtering
type OutgoingFilterMsg struct {
	agent_ip           string
	startFilterPackets bool
}

func NewOutgoingFilterMsg(agent_ip string, startFilterPackets bool) OutgoingFilterMsg {
	return OutgoingFilterMsg{
		agent_ip:           agent_ip,
		startFilterPackets: startFilterPackets}
}

func (filterMsg OutgoingFilterMsg) String() string {
	return "type: Filter; payload: Start / stop filtering packets"
}

func (filterMsg OutgoingFilterMsg) GetType() string {
	return "Filter"
}
