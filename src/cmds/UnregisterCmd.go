package cmds

import "subsystems"
import "visitors"

type UnregisterCmd struct {
	agentRegistry *subsystems.AgentRegistry
	agent_ip      string
	handler_ip    string
	handler_port  string
	agent_port    string
}

func NewUnregisterCmd(aIP string, hIP string, hPort string, aPort string) UnregisterCmd {
	return UnregisterCmd{
		agentRegistry: visitors.AgentRegVisitor.AgentReg,
		agent_ip:      aIP,
		handler_ip:    hIP,
		handler_port:  hPort,
		agent_port:    aPort}
}

func (unregCmd UnregisterCmd) ExecCmd() {
	unregCmd.agentRegistry.RemoveAgent(unregCmd.agent_ip, unregCmd.agent_port)
}
