package outgoingMsg

import "encoding/json"

type OutgoingMsg interface {
	GetType() string
	String() string
}

type OutgoingMsgBuilder struct {
	MsgType string
	Payload interface{}
}

func EncodeMsg(msg OutgoingMsg) ([]byte, error) {
	builder := OutgoingMsgBuilder{MsgType: msg.GetType(), Payload: msg}

	builderData, err := json.Marshal(builder)

	return builderData, err
}
