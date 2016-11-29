package msgs

import "network"

type RegisterMsg struct {
	record network.AgentRecord
}

func newRegisterMsg(rec network.AgentRecord) RegisterMsg {
	return RegisterMsg{
		record: rec}
}

func (regMsg RegisterMsg) String() string {
	return "type: register; payload: agent= " + regMsg.record.GetAgHostname()
}

func (regMsg RegisterMsg) GetType() string {
	return "register"
}
