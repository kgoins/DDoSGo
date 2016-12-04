package outgoingMsg

type OutgoingDebug struct {
	MsgText string
	Id      int
}

func NewOutgoingDebug(msg string) OutgoingDebug {
	return OutgoingDebug{MsgText: msg, Id: 12}
}

func (debug OutgoingDebug) String() string {
	return "type: Debug; payload: " + debug.MsgText + string(debug.Id)
}

func (debug OutgoingDebug) GetType() string {
	return "Debug"
}
