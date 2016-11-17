package msgs

type DebugMsg struct {
	MsgText string
	Id      int
}

func NewDebugMsg(msg string) DebugMsg {
	return DebugMsg{MsgText: msg, Id: 12}
}

func (debugMsg DebugMsg) String() string {
	return "type: debug; payload: " + debugMsg.MsgText + string(debugMsg.Id)
}

func (debugMsg DebugMsg) GetType() string {
	return "debug"
}
