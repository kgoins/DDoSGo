package msgs

import "cmds"

// Filter Message Function, Flag to Start Filtering
type StopEnforcerMsg struct {
	agent_ip   string
	agent_port string
}

func NewStopEnforcerMsg(agent_ip string, agent_port string) StopEnforcerMsg {
	return StopEnforcerMsg{
		agent_ip:   agent_ip,
		agent_port: agent_port,
	}
}

func (stopEnforcerMsg StopEnforcerMsg) String() string {
	return "type: StopEnforcer; payload: Stop Filtering Packets"
}

func (stopEnforcerMsg StopEnforcerMsg) GetType() string {
	return "StopEnforcer"
}

func (stopEnforcerMsg StopEnforcerMsg) BuildCmd() cmds.Cmd {
	return cmds.NewStopEnforcerCmd(stopEnforcerMsg.agent_ip, stopEnforcerMsg.agent_port)
}
