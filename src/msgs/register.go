package msgs

import "cmds"

type RegisterMsg struct {
	Agent_ip     string
	Handler_ip   string
	Handler_port string
	Agent_port   string
	Traceroute   []string
}

func NewRegisterMsg(aIP string, hIP string, port string, trace []string, aPort string) RegisterMsg {
	return RegisterMsg{
		Agent_ip:     aIP,
		Handler_ip:   hIP,
		Handler_port: port,
		Agent_port:   aPort,
		Traceroute:   trace}
}

func (regMsg RegisterMsg) String() string {
	return "type: Register; payload: agent= " + regMsg.Agent_ip + regMsg.Agent_port
}

func (regMsg RegisterMsg) GetType() string {
	return "Register"
}

func (regMsg RegisterMsg) BuildCmd() cmds.Cmd {
	return cmds.NewRegisterCmd(regMsg.Agent_ip, regMsg.Handler_ip, regMsg.Handler_port, regMsg.Traceroute, regMsg.Agent_port)
}
