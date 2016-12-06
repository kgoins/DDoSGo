package msgs

import "cmds"

type UnregisterMsg struct {
	Agent_ip     string
	Handler_ip   string
	Handler_port string
	Agent_port   string
}

func NewUnregisterMsg(aIP string, hIP string, port string, aPort string) UnregisterMsg {
	return UnregisterMsg{
		Agent_ip:     aIP,
		Handler_ip:   hIP,
		Handler_port: port,
		Agent_port:   aPort}
}

func (unregMsg UnregisterMsg) String() string {
	return "type: Unregister; payload: agent= " + unregMsg.Agent_ip + unregMsg.Agent_port
}

func (unregMsg UnregisterMsg) GetType() string {
	return "Unregister"
}

func (unregMsg UnregisterMsg) BuildCmd() cmds.Cmd {
	return cmds.NewUnregisterCmd(unregMsg.Agent_ip, unregMsg.Handler_ip, unregMsg.Handler_port, unregMsg.Agent_port)
}
