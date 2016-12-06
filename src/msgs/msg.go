package msgs

import "encoding/json"
import "fmt"

import "cmds"

type Msg interface {
	BuildCmd() cmds.Cmd
	String() string
	GetType() string
}

type MsgBuilder struct {
	MsgType string
	Payload interface{}
}

func BuildMsg(msgBytes []byte) Msg {
	var rawMsg json.RawMessage
	builder := MsgBuilder{
		Payload: &rawMsg,
	}

	if err := json.Unmarshal(msgBytes, &builder); err != nil {
		fmt.Println(err)
	}

	switch builder.MsgType {
	case "Debug":
		msg := Debug{}
		json.Unmarshal(rawMsg, &msg)
		return msg

	case "DataStream":
		msg := DataStream{}
		json.Unmarshal(rawMsg, &msg)
		return msg

	case "Register":
		msg := RegisterMsg{}
		json.Unmarshal(rawMsg, &msg)
		return msg

	case "Filter":
		msg := FilterMsg{}
		json.Unmarshal(rawMsg, &msg)
		return msg

	default:
		fmt.Println("unknown message type: %s", builder.MsgType)
		return nil
	}
}
