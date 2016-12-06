package outgoingMsg

import "strconv"

type OutgoingRegisterMsg struct {
	Agent_ip     string
	Handler_ip   string
	Handler_port string
	Agent_port   string
	Traceroute   []string
	RemoveFlag   bool
}

func NewOutgoingRegisterMsg(aIP string, hIP string, port string, trace []string, aPort string, addRem bool) OutgoingRegisterMsg {
	return OutgoingRegisterMsg{
		Agent_ip:     aIP,
		Handler_ip:   hIP,
		Handler_port: port,
		Agent_port:   aPort,
		Traceroute:   trace,
		RemoveFlag:   addRem}
}

func (regMsg OutgoingRegisterMsg) String() string {
	return "type: Register - " + strconv.FormatBool(regMsg.RemoveFlag) + "; payload: agent= " + regMsg.Agent_ip + regMsg.Agent_port
}

func (regMsg OutgoingRegisterMsg) GetType() string {
	return "Register"
}
