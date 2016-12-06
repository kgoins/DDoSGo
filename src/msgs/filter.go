package msgs

import "cmds"

// Filter Message Function, Flag to Start Filtering
type FilterMsg struct {
	agent_ip   string
	agent_port string
}

func NewFilterMsg(agent_ip string, agent_port string) FilterMsg {
	return FilterMsg{
		agent_ip:   agent_ip,
		agent_port: agent_port,
	}
}

func (filterMsg FilterMsg) String() string {
	return "type: Filter; payload: Start Filtering Packets"
}

func (filterMsg FilterMsg) GetType() string {
	return "Filter"
}

func (filterMsg FilterMsg) BuildCmd() cmds.Cmd {
	return cmds.NewFilterCmd(filterMsg.agent_ip, filterMsg.agent_port)
}
