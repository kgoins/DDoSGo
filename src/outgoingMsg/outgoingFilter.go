package outgoingMsg

// Filter Message Function, Flag to Start Filtering
type OutgoingFilterMsg struct {
	agent_ip   string
	agent_port string
}

func NewOutgoingFilterMsg(agent_ip string, agent_port string) OutgoingFilterMsg {
	return OutgoingFilterMsg{
		agent_ip:   agent_ip,
		agent_port: agent_port}
}

func (filterMsg OutgoingFilterMsg) String() string {
	return "type: Filter; payload: Start Filtering Packets"
}

func (filterMsg OutgoingFilterMsg) GetType() string {
	return "Filter"
}
