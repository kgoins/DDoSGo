package msgs

import "cmds"

type Debug struct {
	MsgText string
	Id      int
}

func NewDebug(msg string) Debug {
	return Debug{MsgText: msg, Id: 12}
}

func (debug Debug) String() string {
	return "type: Debug; payload: " + debug.MsgText + string(debug.Id)
}

func (debug Debug) GetType() string {
	return "Debug"
}

func (debug Debug) BuildCmd() cmds.Cmd {
	command := cmds.NewDebug(debug.MsgText, debug.Id)
	return command
}
