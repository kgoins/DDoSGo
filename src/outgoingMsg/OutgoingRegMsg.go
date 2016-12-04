package outgoingMsg

type OutgoingRegisterMsg struct {
	Agent_ip      string
	Handler_ip    string
	Handler_port  string
	Traceroute    []string
	AddRemoveFlag bool
}

func NewOutgoingRegisterMsg(aIP string, hIP string, port string, trace []string, addRem bool) OutgoingRegisterMsg {
	return OutgoingRegisterMsg{
		Agent_ip:      aIP,
		Handler_ip:    hIP,
		Handler_port:  port,
		Traceroute:    trace,
		AddRemoveFlag: addRem}
}

func (regMsg OutgoingRegisterMsg) String() string {
	return "type: Register; payload: agent= " + regMsg.Agent_ip
}

func (regMsg OutgoingRegisterMsg) GetType() string {
	return "Register"
}
