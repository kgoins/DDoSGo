package msgs

// import "Cmd"
import "encoding/json"
import "log"

type Msg interface {
	// BuildCommand()
	String() string
	GetType() string
}

type MsgBuilder struct {
	MsgType string
	Payload interface{}
}

func EncodeMsg(msg Msg) []byte {
	// msgData, _ := json.Marshal(msg)
	builder := MsgBuilder{MsgType: msg.GetType(), Payload: msg}

	builderData, _ := json.Marshal(builder)

	return builderData
}

func BulidMsg(msgBytes []byte) Msg {
	var rawMsg json.RawMessage
	builder := MsgBuilder{
		Payload: &rawMsg,
	}

	if err := json.Unmarshal(msgBytes, &builder); err != nil {
		log.Fatal(err)
	}

	switch builder.MsgType {
	case "debug":
		msg := DebugMsg{}
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Fatal(err)
			return nil
		} else {
			return msg
		}
	default:
		log.Fatalf("unknown message type: %q", builder.MsgType)
		return nil
	}
}
