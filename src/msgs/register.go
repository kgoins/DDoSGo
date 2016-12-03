package msgs

type RegisterMsg struct {
	Agent_ip     string
	Handler_ip   string
	Handler_port string
	Traceroute   []string
}

func NewRegisterMsg(aIP string, hIP string, port string, trace []string) RegisterMsg {
	return RegisterMsg{
		Agent_ip:     aIP,
		Handler_ip:   hIP,
		Handler_port: port,
		Traceroute:   trace}
}

func (regMsg RegisterMsg) String() string {
	return "type: register; payload: agent= " + regMsg.Agent_ip
}

func (regMsg RegisterMsg) GetType() string {
	return "Register"
}
