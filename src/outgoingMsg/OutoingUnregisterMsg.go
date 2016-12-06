package outgoingMsg

type OutgoingUnregisterMsg struct {
	Agent_ip     string
	Handler_ip   string
	Handler_port string
	Agent_port   string
}

func NewOutgoingUnegisterMsg(aIP string, hIP string, port string, aPort string) OutgoingUnregisterMsg {
	return OutgoingUnregisterMsg{
		Agent_ip:     aIP,
		Handler_ip:   hIP,
		Handler_port: port,
		Agent_port:   aPort}
}

func (unregMsg OutgoingUnregisterMsg) String() string {
	return "type: Unregister" + "; payload: agent= " + unregMsg.Agent_ip + unregMsg.Agent_port
}

func (unregMsg OutgoingUnregisterMsg) GetType() string {
	return "Unregister"
}
