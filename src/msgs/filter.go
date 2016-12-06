package msgs

import "cmds"

// Filter Message Function, Flag to Start Filtering
type FilterMsg struct {
	agent_ip           string
	startFilterPackets bool
}

func NewFilterMsg(agent_ip string, startFilterPackets bool) FilterMsg {
	return FilterMsg{
		agent_ip:           agent_ip,
		startFilterPackets: startFilterPackets,
	}
}

func (filterMsg FilterMsg) String() string {
	return "type: Filter; payload: Start / stop filtering packets"
}

func (filterMsg FilterMsg) GetType() string {
	return "Filter"
}

func (filterMsg FilterMsg) BuildCmd() cmds.Cmd {
	return cmds.NewFilterCmd(filterMsg.agent_ip, filterMsg.startFilterPackets)
}
