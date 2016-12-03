package msgs

type RegisterMsg struct {
	agent_host   string
	handler_host string
	handler_port int
	traceroute   []string
}

func NewRegisterMsg(aHost string, hHost string, port int, trace []string) RegisterMsg {
	return RegisterMsg{
		agent_host:   aHost,
		handler_host: hHost,
		handler_port: port,
		traceroute:   trace}
}

func (regMsg RegisterMsg) String() string {
	return "type: register; payload: agent= " + regMsg.agent_host
}

func (regMsg RegisterMsg) GetType() string {
	return "register"
}
